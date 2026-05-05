package servicehandler

import (
	"context"
	"errors"
	"fmt"
	"server/model"
	pb "server/protobuf/gen"
	"strings"
	"time"

	"github.com/Jinchenyuan/wego"
	"github.com/Jinchenyuan/wego/logger"
)

type Storyline struct {
	log *logger.Logger
}

func NewStoryline(log *logger.Logger) *Storyline {
	return &Storyline{log: resolveStorylineLogger(log)}
}

func resolveStorylineLogger(log *logger.Logger) *logger.Logger {
	if log != nil {
		return log
	}
	if globalLog := wego.GetGlobalLogger(); globalLog != nil {
		return globalLog
	}
	return logger.GetLogger("storyline.service")
}

func (s *Storyline) ListStorylines(ctx context.Context, req *pb.ListStorylinesRequest, rsp *pb.ListStorylinesResponse) error {
	m := wego.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}

	rows, err := model.ListStorylineSummaries(ctx, m.DB, req.GetAccountId())
	if err != nil {
		return err
	}

	rsp.Storylines = make([]*pb.StorylineSummary, 0, len(rows))
	for _, row := range rows {
		summary := &pb.StorylineSummary{
			Id:            row.ID,
			Title:         row.Title,
			Description:   row.Description,
			CoverPhotoUri: nullStringValue(row.CoverPhotoURI),
			NodeCount:     row.NodeCount,
			UpdatedAt:     row.UpdatedAt.Format(time.RFC3339Nano),
		}
		if row.LatestNodeTitle.Valid {
			summary.LatestNodeTitle = row.LatestNodeTitle.String
		}
		if row.LatestNodeDate.Valid {
			summary.LatestNodeDate = row.LatestNodeDate.Time.Format(time.RFC3339Nano)
		}
		rsp.Storylines = append(rsp.Storylines, summary)
	}
	return nil
}

func (s *Storyline) GetStoryline(ctx context.Context, req *pb.GetStorylineRequest, rsp *pb.StorylineDetail) error {
	m := wego.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}
	storyline, err := model.FindStorylineByID(ctx, m.DB, req.GetAccountId(), req.GetId())
	if err != nil {
		return err
	}
	fillDetail(rsp, storyline)
	return nil
}

func (s *Storyline) CreateStoryline(ctx context.Context, req *pb.CreateStorylineRequest, rsp *pb.StorylineMutationResponse) error {
	if err := validateStorylineMutation(req.GetTitle(), req.GetNodes()); err != nil {
		rsp.Code = 400
		rsp.Message = err.Error()
		return nil
	}

	m := wego.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}
	now := time.Now().UTC()
	storyline := model.NewStorylineRecord(
		req.GetAccountId(),
		strings.TrimSpace(req.GetTitle()),
		req.GetDescription(),
		toNullableString(req.GetCoverPhotoUri()),
		now,
	)
	nodes, err := buildNodes(req.GetNodes(), storyline.ID, now)
	if err != nil {
		rsp.Code = 400
		rsp.Message = err.Error()
		return nil
	}

	if err := model.CreateStoryline(ctx, m.DB, storyline, nodes); err != nil {
		return err
	}
	storyline.Nodes = nodes
	rsp.Success = true
	rsp.Code = 201
	rsp.Message = "Storyline created"
	rsp.Storyline = toProtoDetail(storyline)
	return nil
}

func (s *Storyline) UpdateStoryline(ctx context.Context, req *pb.UpdateStorylineRequest, rsp *pb.StorylineMutationResponse) error {
	if err := validateStorylineMutation(req.GetTitle(), req.GetNodes()); err != nil {
		rsp.Code = 400
		rsp.Message = err.Error()
		return nil
	}

	m := wego.GetGlobalMesa()
	if m == nil {
		return fmt.Errorf("failed to get global mesa")
	}
	existing, err := model.LoadStorylineForUpdate(ctx, m.DB, req.GetAccountId(), req.GetId())
	if err != nil {
		if errors.Is(err, model.ErrStorylineNotFound) {
			return err
		}
		return err
	}

	now := time.Now().UTC()
	existing.Title = strings.TrimSpace(req.GetTitle())
	existing.Description = req.GetDescription()
	existing.CoverPhotoURI = toNullString(req.GetCoverPhotoUri())
	existing.UpdatedAt = now

	nodes, err := buildNodes(req.GetNodes(), existing.ID, now)
	if err != nil {
		rsp.Code = 400
		rsp.Message = err.Error()
		return nil
	}
	if err := model.ReplaceStorylineNodes(ctx, m.DB, existing, nodes); err != nil {
		return err
	}
	existing.Nodes = nodes

	rsp.Success = true
	rsp.Code = 200
	rsp.Message = "Storyline updated"
	rsp.Storyline = toProtoDetail(existing)
	return nil
}
