package tuic

import (
	"context"
	"time"

	"github.com/sagernet/quic-go"
	"github.com/sagernet/quic-go/congestion"
	congestion_meta1 "github.com/sagernet/sing-quic/congestion_meta1"
	congestion_meta2 "github.com/sagernet/sing-quic/congestion_meta2"
	"github.com/sagernet/sing/common/ntp"
)

func setCongestion(ctx context.Context, connection quic.Connection, congestionName string) {
	timeFunc := ntp.TimeFuncFromContext(ctx)
	if timeFunc == nil {
		timeFunc = time.Now
	}
	switch congestionName {
	case "cubic":
		connection.SetCongestionControl(
			congestion_meta1.NewCubicSender(
				congestion_meta1.DefaultClock{TimeFunc: timeFunc},
				congestion_meta1.GetInitialPacketSize(connection.RemoteAddr()),
				false,
				nil,
			),
		)
	case "new_reno":
		connection.SetCongestionControl(
			congestion_meta1.NewCubicSender(
				congestion_meta1.DefaultClock{TimeFunc: timeFunc},
				congestion_meta1.GetInitialPacketSize(connection.RemoteAddr()),
				true,
				nil,
			),
		)
	case "bbr_meta_v1":
		connection.SetCongestionControl(congestion_meta1.NewBBRSender(
			congestion_meta1.DefaultClock{TimeFunc: timeFunc},
			congestion_meta1.GetInitialPacketSize(connection.RemoteAddr()),
			congestion_meta1.InitialCongestionWindow*congestion_meta1.InitialMaxDatagramSize,
			congestion_meta1.DefaultBBRMaxCongestionWindow*congestion_meta1.InitialMaxDatagramSize,
		))
	case "bbr":
		connection.SetCongestionControl(congestion_meta2.NewBbrSender(
			congestion_meta2.DefaultClock{TimeFunc: timeFunc},
			congestion_meta2.GetInitialPacketSize(connection.RemoteAddr()),
			congestion.ByteCount(congestion_meta1.InitialCongestionWindow),
		))
	}
}
