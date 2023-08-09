package base

import (
	"context"
	"github.com/capell/capell_scan/proto/base"
	"github.com/capell/capell_scan/rpc"
)

func BaseLp(ctx context.Context, req *base.MsgBaseLp) (resp *base.MsgBaseLpResponse, err error) {
	r, err := rpc.SwapClient.BaseLp(ctx, req)

	if err != nil {
		return
	}

	return r, nil
}
