package stannp

import (
	"context"
	"github.com/copilotiq/stannp-client-golang/address"
	"github.com/copilotiq/stannp-client-golang/letter"
	"github.com/copilotiq/stannp-client-golang/util"
	"io"
	"os"
)

// Client interface is for mocking / testing. Implement it however you wish!
// A standard set of mocks however is available via MockClient
type Client interface {
	GetPDFContents(ctx context.Context, pdfURL string) (*letter.PDFRes, *util.APIError)
	SavePDFContents(pdfContents io.Reader) (*os.File, *util.APIError)
	SendLetter(ctx context.Context, req *letter.SendReq) (*letter.SendRes, *util.APIError)
	ValidateAddress(ctx context.Context, req *address.ValidateReq) (*address.ValidateRes, *util.APIError)
}
