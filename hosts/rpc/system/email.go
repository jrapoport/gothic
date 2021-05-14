package system

import (
	"context"

	"github.com/google/uuid"
	"github.com/jrapoport/gothic/api/grpc/rpc/system"
	"github.com/jrapoport/gothic/hosts/rpc"
	"github.com/jrapoport/gothic/mail"
	"google.golang.org/grpc/codes"
)

func (s *systemServer) SendEmail(ctx context.Context, req *system.EmailRequest) (*system.EmailResponse, error) {
	uid, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, s.RPCError(codes.InvalidArgument, err)
	}
	sub := req.GetSubject()
	content := mail.Content{
		Type:      0,
		Body:      "",
		Plaintext: "",
	}
	switch req.GetContent().(type) {
	case *system.EmailRequest_Body:
		content.Type = mail.Template
		content.Body = req.GetBody()
	case *system.EmailRequest_Markdown:
		content.Type = mail.Markdown
		content.Body = req.GetMarkdown()
	default:
		content.Type = mail.HTML
		content.Body = req.GetHtml()
		content.Plaintext = req.GetPlaintext()
	}
	rtx := rpc.RequestContext(ctx)
	sent, err := s.API.SendEmail(rtx, uid, sub, content)
	if err != nil {
		return nil, s.RPCError(codes.Internal, err)
	}
	res := &system.EmailResponse{
		Sent: sent,
	}
	return res, nil
}
