package websocket

import (
	"context"
	"time"
)

// StartHearthBeat keeps sending keep-alive requests at a specified interval.
// If no interval is specified, a default is used.
// It's recommended to use the default value unless you have a good reason not to do so.
//
// See Websocket.SendHeartBeat for more information.
func (w *Websocket) StartHearthBeat(ctx context.Context, duration ...time.Duration) {
	d := time.Millisecond * 49000
	if len(duration) > 0 {
		d = duration[0]
	}

	ticker := time.NewTicker(d)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			w.SendHeartBeat()
		}
	}
}
