package ginhandler

import (
	"errors"
	pb "server/protobuf/gen"
	"strings"
)

func toCreateRequest(accountID uint32, req CreateStorylineRequest) *pb.CreateStorylineRequest {
	return &pb.CreateStorylineRequest{
		AccountId:     accountID,
		Title:         req.Title,
		Description:   req.Description,
		CoverPhotoUri: normalizeOptionalString(req.CoverPhotoURI),
		Nodes:         toProtoNodeInputs(req.Nodes),
	}
}

func toUpdateRequest(accountID uint32, req UpdateStorylineRequest) *pb.UpdateStorylineRequest {
	return &pb.UpdateStorylineRequest{
		AccountId:     accountID,
		Id:            req.ID,
		Title:         req.Title,
		Description:   req.Description,
		CoverPhotoUri: normalizeOptionalString(req.CoverPhotoURI),
		Nodes:         toProtoNodeInputs(req.Nodes),
	}
}

func toProtoNodeInputs(inputs []StorylineNodeInput) []*pb.StorylineNodeInput {
	ret := make([]*pb.StorylineNodeInput, 0, len(inputs))
	for _, input := range inputs {
		ret = append(ret, &pb.StorylineNodeInput{
			Id:        normalizeOptionalString(input.ID),
			Title:     input.Title,
			Date:      input.Date,
			Note:      input.Note,
			Location:  input.Location,
			PhotoUri:  normalizeOptionalString(input.PhotoURI),
			SortOrder: input.SortOrder,
		})
	}
	return ret
}

func toListResponse(rsp *pb.ListStorylinesResponse) ListStorylinesResponse {
	items := make([]StorylineSummary, 0, len(rsp.GetStorylines()))
	for _, item := range rsp.GetStorylines() {
		items = append(items, StorylineSummary{
			ID:              item.GetId(),
			Title:           item.GetTitle(),
			Description:     item.GetDescription(),
			CoverPhotoURI:   toOptionalString(item.GetCoverPhotoUri()),
			NodeCount:       item.GetNodeCount(),
			LatestNodeTitle: toOptionalString(item.GetLatestNodeTitle()),
			LatestNodeDate:  toOptionalString(item.GetLatestNodeDate()),
			UpdatedAt:       item.GetUpdatedAt(),
		})
	}
	return ListStorylinesResponse{Storylines: items}
}

func toMutationResponse(rsp *pb.StorylineMutationResponse) StorylineMutationResponse {
	return StorylineMutationResponse{
		Success:   rsp.GetSuccess(),
		Storyline: toDetail(rsp.GetStoryline()),
		Message:   toOptionalString(rsp.GetMessage()),
	}
}

func toDetail(detail *pb.StorylineDetail) StorylineDetail {
	nodes := make([]StorylineNode, 0, len(detail.GetNodes()))
	for _, node := range detail.GetNodes() {
		nodes = append(nodes, StorylineNode{
			ID:        node.GetId(),
			Title:     node.GetTitle(),
			Date:      node.GetDate(),
			Note:      node.GetNote(),
			Location:  node.GetLocation(),
			PhotoURI:  toOptionalString(node.GetPhotoUri()),
			SortOrder: node.GetSortOrder(),
		})
	}
	return StorylineDetail{
		ID:            detail.GetId(),
		Title:         detail.GetTitle(),
		Description:   detail.GetDescription(),
		CoverPhotoURI: toOptionalString(detail.GetCoverPhotoUri()),
		CreatedAt:     detail.GetCreatedAt(),
		UpdatedAt:     detail.GetUpdatedAt(),
		Nodes:         nodes,
	}
}

func validateNodes(nodes []StorylineNodeInput) error {
	if len(nodes) == 0 {
		return errors.New("storyline must contain at least one node")
	}
	for _, node := range nodes {
		if strings.TrimSpace(node.Title) == "" {
			return errors.New("storyline node title is required")
		}
		if strings.TrimSpace(node.Date) == "" {
			return errors.New("storyline node date is required")
		}
	}
	return nil
}

func normalizeOptionalString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func toOptionalString(value string) *string {
	if value == "" {
		return nil
	}
	ret := value
	return &ret
}
