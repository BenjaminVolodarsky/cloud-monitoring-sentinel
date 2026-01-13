package vm

import "context"

type QueryOptions struct {
	Expr  string
	Start string
	End   string
	Step  string
}

func (c *Client) Query(ctx context.Context, opts QueryOptions) ([]byte, error) {
	params := map[string]string{
		"query": opts.Expr,
	}

	if opts.Start != "" {
		params["start"] = opts.Start
	}
	if opts.End != "" {
		params["end"] = opts.End
	}
	if opts.Step != "" {
		params["step"] = opts.Step
	}

	return c.doGET(
		ctx,
		"/select/0/prometheus/api/v1/query",
		params,
	)
}
