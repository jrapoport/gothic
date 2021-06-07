package admin

import (
	"context"

	api "github.com/jrapoport/gothic/api/grpc/rpc"
	"github.com/jrapoport/gothic/api/grpc/rpc/admin"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/models/auditlog"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *adminServer) SearchAuditLogs(ctx context.Context,
	req *api.SearchRequest) (*admin.AuditLogsResult, error) {
	if req == nil {
		return nil, s.RPCError(codes.InvalidArgument, nil)
	}
	rtx, err := s.adminRequestContext(ctx)
	if err != nil {
		return nil, s.RPCError(codes.PermissionDenied, err)
	}
	filters := req.Filters.AsMap()
	page := rpc.PaginateRequest(req)
	s.Debugf("search audit logs: %v", req)
	logs, err := s.API.SearchAuditLogs(rtx, filters, page)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	s.Debugf("found %d logs", len(logs))
	auditLogs := make([]*admin.AuditLog, len(logs))
	for i, log := range logs {
		auditLogs[i] = auditLogResponse(log)
	}
	res := &admin.AuditLogsResult{
		Logs: auditLogs,
		Page: rpc.PaginateResponse(page),
	}
	return res, nil
}

func auditLogResponse(log *auditlog.AuditLog) *admin.AuditLog {
	res := &admin.AuditLog{
		Id:        uint64(log.ID),
		Type:      admin.AuditLog_Type(log.Type),
		Action:    log.Action.String(),
		UserId:    log.UserID.String(),
		CreatedAt: timestamppb.New(log.CreatedAt),
	}
	res.Fields, _ = structpb.NewStruct(log.Fields)
	return res
}
