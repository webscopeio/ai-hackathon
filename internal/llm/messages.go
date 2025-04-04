package llm

import (
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"golang.org/x/net/context"
)

func (c *Client) NewMessage(
	ctx context.Context,
	body anthropic.MessageNewParams,
	opts ...option.RequestOption,
) (res *anthropic.Message, err error) {
	return c.client.Messages.New(ctx, body, opts...)
}
