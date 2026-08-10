package service

import (
	"context"
	"fmt"
)

type BatchOperationService struct {
	svc Service
}

func NewBatchOperationService(svc Service) *BatchOperationService {
	return &BatchOperationService{svc: svc}
}

type BatchResult struct {
	Succeeded int      `json:"succeeded"`
	Failed    int      `json:"failed"`
	Errors    []string `json:"errors,omitempty"`
}

func (s *BatchOperationService) BatchSubmit(ctx context.Context, ids []string, userID string) (*BatchResult, error) {
	r := &BatchResult{}
	for _, id := range ids {
		if err := s.svc.SubmitForReview(ctx, id, userID); err != nil {
			r.Failed++
			r.Errors = append(r.Errors, fmt.Sprintf("%s: %v", id, err))
		} else {
			r.Succeeded++
		}
	}
	return r, nil
}

func (s *BatchOperationService) BatchApprove(ctx context.Context, ids []string, approverID string) (*BatchResult, error) {
	r := &BatchResult{}
	for _, id := range ids {
		if err := s.svc.ApproveEntry(ctx, id, approverID); err != nil {
			r.Failed++
			r.Errors = append(r.Errors, fmt.Sprintf("%s: %v", id, err))
		} else {
			r.Succeeded++
		}
	}
	return r, nil
}

func (s *BatchOperationService) BatchPost(ctx context.Context, ids []string) (*BatchResult, error) {
	r := &BatchResult{}
	for _, id := range ids {
		if err := s.svc.PostEntry(ctx, id); err != nil {
			r.Failed++
			r.Errors = append(r.Errors, fmt.Sprintf("%s: %v", id, err))
		} else {
			r.Succeeded++
		}
	}
	return r, nil
}

func (s *BatchOperationService) BatchCancel(ctx context.Context, ids []string) (*BatchResult, error) {
	r := &BatchResult{}
	for _, id := range ids {
		if err := s.svc.CancelEntry(ctx, id); err != nil {
			r.Failed++
			r.Errors = append(r.Errors, fmt.Sprintf("%s: %v", id, err))
		} else {
			r.Succeeded++
		}
	}
	return r, nil
}
