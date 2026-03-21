package commands

import (
	"context"
	"fmt"
)

// approveCommand returns the command definition for approving shield actions.
func approveCommand() Definition {
	return Definition{
		Name:        "approve",
		Description: "Approve a pending shield approval request",
		Usage:       "approve <approval_id>",
		Handler:     handleApprove,
	}
}

// rejectCommand returns the command definition for rejecting shield actions.
func rejectCommand() Definition {
	return Definition{
		Name:        "reject",
		Description: "Reject a pending shield approval request",
		Usage:       "reject <approval_id>",
		Handler:     handleReject,
	}
}

func handleApprove(ctx context.Context, req Request, rt *Runtime) error {
	approvalID := nthToken(req.Text, 1)
	if approvalID == "" {
		return req.Reply("Usage: approve <approval_id>")
	}

	if rt.ApproveAction == nil {
		return req.Reply("Error: approval system not available")
	}

	if rt.ApproveAction(approvalID) {
		return req.Reply(fmt.Sprintf("✅ Approved request %s", approvalID))
	}
	return req.Reply(fmt.Sprintf("❌ Failed to approve %s (not found or already processed)", approvalID))
}

func handleReject(ctx context.Context, req Request, rt *Runtime) error {
	approvalID := nthToken(req.Text, 1)
	if approvalID == "" {
		return req.Reply("Usage: reject <approval_id>")
	}

	if rt.RejectAction == nil {
		return req.Reply("Error: approval system not available")
	}

	if rt.RejectAction(approvalID) {
		return req.Reply(fmt.Sprintf("🚫 Rejected request %s", approvalID))
	}
	return req.Reply(fmt.Sprintf("❌ Failed to reject %s (not found or already processed)", approvalID))
}
