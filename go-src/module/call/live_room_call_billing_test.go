package call

import (
	"testing"

	callentity "xr-game-server/entity/call"
	userentity "xr-game-server/entity/user"
)

func TestResolveRoomCallPartyIdsByTypes(t *testing.T) {
	const (
		callerId uint64 = 1001
		targetId uint64 = 2002
	)
	tests := []struct {
		name           string
		callerType     uint8
		targetType     uint8
		wantAnchorId   uint64
		wantAudienceId uint64
		wantOK         bool
	}{
		{
			name:           "viewer calls anchor",
			callerType:     userentity.UserTypeNormal,
			targetType:     userentity.UserTypeSeniorAnchor,
			wantAnchorId:   targetId,
			wantAudienceId: callerId,
			wantOK:         true,
		},
		{
			name:           "anchor calls viewer",
			callerType:     userentity.UserTypeBotAnchor,
			targetType:     userentity.UserTypeNormal,
			wantAnchorId:   callerId,
			wantAudienceId: targetId,
			wantOK:         true,
		},
		{
			name:       "two viewers have no room owner",
			callerType: userentity.UserTypeNormal,
			targetType: userentity.UserTypeNormal,
		},
		{
			name:       "two anchors have no ordinary payer",
			callerType: userentity.UserTypeAnchor,
			targetType: userentity.UserTypeSeniorAnchor,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			anchorId, audienceId, ok := resolveRoomCallPartyIdsByTypes(callerId, tt.callerType, targetId, tt.targetType)
			if anchorId != tt.wantAnchorId || audienceId != tt.wantAudienceId || ok != tt.wantOK {
				t.Fatalf("resolveRoomCallPartyIdsByTypes() = (%d, %d, %t), want (%d, %d, %t)", anchorId, audienceId, ok, tt.wantAnchorId, tt.wantAudienceId, tt.wantOK)
			}
		})
	}
}

func TestCallSourceUsesLiveRoom(t *testing.T) {
	tests := []struct {
		name   string
		source uint8
		want   bool
	}{
		{name: "live room", source: callentity.CallOrderSourceLiveRoom, want: true},
		{name: "private message", source: callentity.CallOrderSourcePrivateMessage, want: false},
		{name: "one-to-one room", source: callentity.CallOrderSourceOneToOneRoom, want: true},
		{name: "unknown", source: 99, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := callSourceUsesLiveRoom(tt.source); got != tt.want {
				t.Fatalf("callSourceUsesLiveRoom(%d) = %t, want %t", tt.source, got, tt.want)
			}
		})
	}
}

func TestCallDiamondBillingPartyIdsKeepsLiveRoomSemantics(t *testing.T) {
	const (
		callerId   uint64 = 1001
		receiverId uint64 = 2002
	)
	order := &callentity.CallOrder{
		CallerId:   callerId,
		ReceiverId: receiverId,
		PayerId:    callerId,
		Source:     callentity.CallOrderSourceLiveRoom,
	}
	anchorId, payerId, billable := callDiamondBillingPartyIds(order)
	if anchorId != receiverId || payerId != callerId || !billable {
		t.Fatalf(
			"live-room billing parties = (%d, %d, %t), want (%d, %d, true)",
			anchorId,
			payerId,
			billable,
			receiverId,
			callerId,
		)
	}
}

func TestCallDiamondBillingPartyIdsUsesStoredOneToOnePayer(t *testing.T) {
	const (
		anchorId   uint64 = 1001
		audienceId uint64 = 2002
	)
	tests := []struct {
		name       string
		callerId   uint64
		receiverId uint64
	}{
		{name: "anchor calls audience", callerId: anchorId, receiverId: audienceId},
		{name: "audience calls anchor", callerId: audienceId, receiverId: anchorId},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := &callentity.CallOrder{
				CallerId:   tt.callerId,
				ReceiverId: tt.receiverId,
				PayerId:    audienceId,
				Source:     callentity.CallOrderSourceOneToOneRoom,
			}
			gotAnchorId, gotPayerId, billable := callDiamondBillingPartyIds(order)
			if gotAnchorId != anchorId || gotPayerId != audienceId || !billable {
				t.Fatalf(
					"1v1 billing parties = (%d, %d, %t), want (%d, %d, true)",
					gotAnchorId,
					gotPayerId,
					billable,
					anchorId,
					audienceId,
				)
			}
		})
	}
}

func TestLiveRoomCallTicketPrice(t *testing.T) {
	tests := []struct {
		name    string
		order   *callentity.CallOrder
		enabled bool
		want    float64
	}{
		{
			name:    "source one charges ticket when enabled",
			order:   &callentity.CallOrder{Source: callentity.CallOrderSourceLiveRoom, TicketPrice: 8.5},
			enabled: true,
			want:    8.5,
		},
		{
			name:    "source one skips ticket when disabled",
			order:   &callentity.CallOrder{Source: callentity.CallOrderSourceLiveRoom, TicketPrice: 8.5},
			enabled: false,
		},
		{
			name:    "source three never charges ticket",
			order:   &callentity.CallOrder{Source: callentity.CallOrderSourceOneToOneRoom, TicketPrice: 8.5},
			enabled: true,
		},
		{
			name:    "non-positive ticket is ignored",
			order:   &callentity.CallOrder{Source: callentity.CallOrderSourceLiveRoom, TicketPrice: 0},
			enabled: true,
		},
		{name: "nil order", enabled: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := liveRoomCallTicketPrice(tt.order, tt.enabled); got != tt.want {
				t.Fatalf("liveRoomCallTicketPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
