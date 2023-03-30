package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/falcosecurity/falcosidekick-ui/internal/models"
	"github.com/falcosecurity/falcosidekick-ui/internal/utils"
)

func newQuery(args *models.Arguments) string {
	var filter, priority, rule, source, hostname, tags, since string
	if args.Filter != "" {
		filter = utils.Escape(args.Filter) + "* "
	}
	if args.Priority != "" {
		p := strings.Split(args.Priority, ",")
		for i, j := range p {
			p[i] = fmt.Sprintf("@priority:%v", j)
		}
		priority = fmt.Sprintf("(%v) ", strings.Join(p, " | "))
	}
	if args.Rule != "" {
		r := strings.Split(args.Rule, ",")
		for i, j := range r {
			r[i] = fmt.Sprintf("@rule:'%v'", j)
		}
		rule = fmt.Sprintf("(%v) ", strings.Join(r, " | "))
	}
	if args.Source != "" {
		r := strings.Split(args.Source, ",")
		for i, j := range r {
			r[i] = fmt.Sprintf("@source:%v", j)
		}
		source = fmt.Sprintf("(%v) ", strings.Join(r, " | "))
	}
	if args.Hostname != "" {
		r := strings.Split(utils.Escape(args.Hostname), ",")
		for i, j := range r {
			r[i] = fmt.Sprintf("@hostname:%v", j)
		}
		hostname = fmt.Sprintf("(%v) ", strings.Join(r, " | "))
	}
	if args.Tags != "" {
		r := strings.Split(utils.Escape(args.Tags), ",")
		for i, j := range r {
			r[i] = fmt.Sprintf("@tags:%v", j)
		}
		tags = fmt.Sprintf("(%v) ", strings.Join(r, " | "))
	}
	var s int64
	if t := int64(utils.ConvertToSeconds(args.Since)); t != 0 {
		s = time.Now().UTC().Add(-1 * time.Duration(t) * time.Second).UnixMicro()
	}
	since = fmt.Sprintf("@timestamp:[%v inf]", s)

	return fmt.Sprintf("%v%v%v%v%v%v%v", filter, priority, rule, source, hostname, tags, since)
}
