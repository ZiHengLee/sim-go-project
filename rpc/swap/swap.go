package swap

import (
	"context"
	"github.com/capell/capell_scan/proto/swap"
	"github.com/capell/capell_scan/rpc"
)

func BaseLp(ctx context.Context, req *swap.MsgBaseLp) (resp *swap.MsgBaseLpResponse, err error) {
	r, err := rpc.SwapClient.BaseLp(ctx, req)

	if err != nil {
		return
	}

	return r, nil
}
