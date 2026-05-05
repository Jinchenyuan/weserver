package servicehandler

import (
	"database/sql"
	"errors"
	"server/model"
	pb "server/protobuf/gen"
	"sort"
	"strings"
	"time"
)

func validateStorylineMutation(title string, nodes []*pb.StorylineNodeInput) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("storyline title is required")
	}
	if len(nodes) == 0 {
		return errors.New("storyline must contain at least one node")
	}
	for _, node := range nodes {
		if strings.TrimSpace(node.GetTitle()) == "" {
			return errors.New("storyline node title is required")
		}
		if strings.TrimSpace(node.GetDate()) == "" {
			return errors.New("storyline node date is required")
		}
	}
	return nil
}

func buildNodes(inputs []*pb.StorylineNodeInput, storylineID string, now time.Time) ([]*model.StorylineNode, error) {
	nodeInputs := make([]model.StorylineInput, 0, len(inputs))
	for _, input := range inputs {
		date, err := parseStorylineTime(input.GetDate())
		if err != nil {
			return nil, err
		}
		nodeInputs = append(nodeInputs, model.StorylineInput{
			ID:        input.GetId(),
			Title:     strings.TrimSpace(input.GetTitle()),
			Date:      date.UTC(),
			Note:      input.GetNote(),
			Location:  input.GetLocation(),
			PhotoURI:  toNullableString(input.GetPhotoUri()),
			SortOrder: input.GetSortOrder(),
		})
	}
	nodes := model.NormalizeStorylineNodes(storylineID, now, nodeInputs)
	return nodes, nil
}

func parseStorylineTime(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05",
	}

	var lastErr error
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed.UTC(), nil
		}
		lastErr = err
	}

	return time.Time{}, lastErr
}

func toProtoDetail(storyline *model.Storyline) *pb.StorylineDetail {
	nodes := append([]*model.StorylineNode(nil), storyline.Nodes...)
	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i].SortOrder < nodes[j].SortOrder
	})

	ret := &pb.StorylineDetail{
		Id:            storyline.ID,
		Title:         storyline.Title,
		Description:   storyline.Description,
		CoverPhotoUri: nullStringValue(storyline.CoverPhotoURI),
		CreatedAt:     storyline.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:     storyline.UpdatedAt.Format(time.RFC3339Nano),
		Nodes:         make([]*pb.StorylineNode, 0, len(nodes)),
	}
	for _, node := range nodes {
		ret.Nodes = append(ret.Nodes, &pb.StorylineNode{
			Id:        node.ID,
			Title:     node.Title,
			Date:      node.Date.Format(time.RFC3339Nano),
			Note:      node.Note,
			Location:  node.Location,
			PhotoUri:  nullStringValue(node.PhotoURI),
			SortOrder: node.SortOrder,
		})
	}
	return ret
}

func fillDetail(rsp *pb.StorylineDetail, storyline *model.Storyline) {
	detail := toProtoDetail(storyline)
	*rsp = *detail
}

func toNullableString(value string) *string {
	if value == "" {
		return nil
	}
	ret := value
	return &ret
}

func toNullString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}

func nullStringValue(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
