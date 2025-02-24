// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mac

import (
	"context"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/v3/pkg/band"
	"go.thethings.network/lorawan-stack/v3/pkg/crypto"
	"go.thethings.network/lorawan-stack/v3/pkg/frequencyplans"
	"go.thethings.network/lorawan-stack/v3/pkg/gpstime"
	. "go.thethings.network/lorawan-stack/v3/pkg/networkserver/internal/test"
	"go.thethings.network/lorawan-stack/v3/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/v3/pkg/types"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test"
	"go.thethings.network/lorawan-stack/v3/pkg/util/test/assertions/should"
)

func TestDeviceDefaultChannels(t *testing.T) {
	t.Parallel()

	preBand := band.AU_915_928_RP1_v1_0_2
	au915928Pre102BUplinkChannels := preBand.UplinkChannels
	au915928Pre102BDownlinkChannels := preBand.DownlinkChannels
	au915928Pre102BChannels := make([]*ttnpb.MACParameters_Channel, 0, len(au915928Pre102BUplinkChannels))
	for i, channel := range preBand.UplinkChannels {
		minDataRate, maxDataRate := ttnpb.DataRateIndex_DATA_RATE_0, ttnpb.DataRateIndex_DATA_RATE_3
		if i >= 64 {
			minDataRate, maxDataRate = ttnpb.DataRateIndex_DATA_RATE_4, ttnpb.DataRateIndex_DATA_RATE_4
		}
		au915928Pre102BChannels = append(au915928Pre102BChannels, &ttnpb.MACParameters_Channel{
			UplinkFrequency:   channel.Frequency,
			DownlinkFrequency: au915928Pre102BDownlinkChannels[i%8].Frequency,
			MinDataRateIndex:  minDataRate,
			MaxDataRateIndex:  maxDataRate,
			EnableUplink:      true,
		})
	}

	postBand := band.AU_915_928_RP1_v1_0_2_RevB
	au915928Post102BUplinkChannels := postBand.UplinkChannels
	au915928Post102BDownlinkChannels := postBand.DownlinkChannels
	au915928Post102BChannels := make([]*ttnpb.MACParameters_Channel, 0, len(au915928Post102BUplinkChannels))
	for i, channel := range au915928Post102BUplinkChannels {
		minDataRate, maxDataRate := ttnpb.DataRateIndex_DATA_RATE_0, ttnpb.DataRateIndex_DATA_RATE_5
		if i >= 64 {
			minDataRate, maxDataRate = ttnpb.DataRateIndex_DATA_RATE_6, ttnpb.DataRateIndex_DATA_RATE_6
		}
		au915928Post102BChannels = append(au915928Post102BChannels, &ttnpb.MACParameters_Channel{
			UplinkFrequency:   channel.Frequency,
			DownlinkFrequency: au915928Post102BDownlinkChannels[i%8].Frequency,
			MinDataRateIndex:  minDataRate,
			MaxDataRateIndex:  maxDataRate,
			EnableUplink:      true,
		})
	}

	for _, tc := range []struct {
		Name          string
		Device        *ttnpb.EndDevice
		Band          *band.Band
		FrequencyPlan *frequencyplans.FrequencyPlan
		Channels      []*ttnpb.MACParameters_Channel
	}{
		{
			Name:     "1.0.1/AU_915_928_FSB_2",
			Device:   &ttnpb.EndDevice{},
			Band:     &band.AU_915_928_TS1_v1_0_1,
			Channels: au915928Pre102BChannels,
		},
		{
			Name:     "1.0.2/AU_915_928_FSB_2",
			Device:   &ttnpb.EndDevice{},
			Band:     &band.AU_915_928_RP1_v1_0_2,
			Channels: au915928Pre102BChannels,
		},
		{
			Name:     "1.0.2B/AU_915_928_FSB_2",
			Device:   &ttnpb.EndDevice{},
			Band:     &band.AU_915_928_RP1_v1_0_2_RevB,
			Channels: au915928Post102BChannels,
		},
	} {
		tc := tc
		test.RunSubtest(t, test.SubtestConfig{
			Name:     tc.Name,
			Parallel: true,
			Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
				t.Helper()

				pb := ttnpb.Clone(tc.Device)

				deviceDefaultChannels := DeviceDefaultChannels(pb, tc.Band, &ttnpb.MACSettings{})
				a.So(deviceDefaultChannels, should.Resemble, tc.Channels)
				a.So(pb, should.Resemble, tc.Device)
			},
		})
	}
}

func TestDeviceDesiredChannels(t *testing.T) {
	t.Parallel()

	preBand := band.AU_915_928_RP1_v1_0_2
	au915928Pre102BUplinkChannels := preBand.UplinkChannels
	au915928Pre102BDownlinkChannels := preBand.DownlinkChannels
	au915928Pre102BChannels := make([]*ttnpb.MACParameters_Channel, 0, len(au915928Pre102BUplinkChannels))
	for i, channel := range preBand.UplinkChannels {
		minDataRate, maxDataRate := ttnpb.DataRateIndex_DATA_RATE_0, ttnpb.DataRateIndex_DATA_RATE_3
		if i >= 64 {
			minDataRate, maxDataRate = ttnpb.DataRateIndex_DATA_RATE_4, ttnpb.DataRateIndex_DATA_RATE_4
		}
		au915928Pre102BChannels = append(au915928Pre102BChannels, &ttnpb.MACParameters_Channel{
			UplinkFrequency:   channel.Frequency,
			DownlinkFrequency: au915928Pre102BDownlinkChannels[i%8].Frequency,
			MinDataRateIndex:  minDataRate,
			MaxDataRateIndex:  maxDataRate,
			EnableUplink:      i >= 8 && i < 16 || i == 65, // FSB2
		})
	}

	postBand := band.AU_915_928_RP1_v1_0_2_RevB
	au915928Post102BUplinkChannels := postBand.UplinkChannels
	au915928Post102BDownlinkChannels := postBand.DownlinkChannels
	au915928Post102BChannels := make([]*ttnpb.MACParameters_Channel, 0, len(au915928Post102BUplinkChannels))
	for i, channel := range postBand.UplinkChannels {
		minDataRate, maxDataRate := ttnpb.DataRateIndex_DATA_RATE_0, ttnpb.DataRateIndex_DATA_RATE_5
		if i >= 64 {
			minDataRate, maxDataRate = ttnpb.DataRateIndex_DATA_RATE_6, ttnpb.DataRateIndex_DATA_RATE_6
		}
		au915928Post102BChannels = append(au915928Post102BChannels, &ttnpb.MACParameters_Channel{
			UplinkFrequency:   channel.Frequency,
			DownlinkFrequency: au915928Post102BDownlinkChannels[i%8].Frequency,
			MinDataRateIndex:  minDataRate,
			MaxDataRateIndex:  maxDataRate,
			EnableUplink:      i >= 8 && i < 16 || i == 65, // FSB2
		})
	}

	for _, tc := range []struct {
		Name          string
		Device        *ttnpb.EndDevice
		Band          *band.Band
		FrequencyPlan *frequencyplans.FrequencyPlan
		Channels      []*ttnpb.MACParameters_Channel
	}{
		{
			Name:          "1.0.1/AU_915_928_FSB_2",
			Device:        &ttnpb.EndDevice{},
			Band:          &band.AU_915_928_TS1_v1_0_1,
			FrequencyPlan: test.FrequencyPlan(test.AUFrequencyPlanID),
			Channels:      au915928Pre102BChannels,
		},
		{
			Name:          "1.0.2/AU_915_928_FSB_2",
			Device:        &ttnpb.EndDevice{},
			Band:          &band.AU_915_928_RP1_v1_0_2,
			FrequencyPlan: test.FrequencyPlan(test.AUFrequencyPlanID),
			Channels:      au915928Pre102BChannels,
		},
		{
			Name:          "1.0.2B/AU_915_928_FSB_2",
			Device:        &ttnpb.EndDevice{},
			Band:          &band.AU_915_928_RP1_v1_0_2_RevB,
			FrequencyPlan: test.FrequencyPlan(test.AUFrequencyPlanID),
			Channels:      au915928Post102BChannels,
		},
	} {
		tc := tc
		test.RunSubtest(t, test.SubtestConfig{
			Name:     tc.Name,
			Parallel: true,
			Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
				t.Helper()

				pb := ttnpb.Clone(tc.Device)

				deviceDesiredChannels := DeviceDesiredChannels(pb, tc.Band, tc.FrequencyPlan, &ttnpb.MACSettings{})
				a.So(deviceDesiredChannels, should.Resemble, tc.Channels)
				a.So(pb, should.Resemble, tc.Device)
			},
		})
	}
}

func TestNewState(t *testing.T) {
	for _, tc := range []struct {
		Name               string
		Device             *ttnpb.EndDevice
		MACState           *ttnpb.MACState
		FrequencyPlanStore *frequencyplans.Store
		ErrorAssertion     func(*testing.T, error) bool
	}{
		{
			Name: "1.0.2/EU868",
			Device: &ttnpb.EndDevice{
				FrequencyPlanId:   test.EUFrequencyPlanID,
				LorawanVersion:    ttnpb.MACVersion_MAC_V1_0_2,
				LorawanPhyVersion: ttnpb.PHYVersion_RP001_V1_0_2_REV_B,
				MacSettings: &ttnpb.MACSettings{
					DesiredRx1Delay: &ttnpb.RxDelayValue{
						Value: ttnpb.RxDelay_RX_DELAY_13,
					},
				},
			},
			MACState: func() *ttnpb.MACState {
				macState := MakeDefaultEU868MACState(ttnpb.Class_CLASS_A, ttnpb.MACVersion_MAC_V1_0_2, ttnpb.PHYVersion_RP001_V1_0_2_REV_B)
				macState.DesiredParameters.Rx1Delay = ttnpb.RxDelay_RX_DELAY_13
				return macState
			}(),
			FrequencyPlanStore: frequencyplans.NewStore(test.FrequencyPlansFetcher),
		},
		{
			Name: "1.0.2/EU868/DesiredMaxEirp",
			Device: &ttnpb.EndDevice{
				FrequencyPlanId:   test.EUFrequencyPlanID,
				LorawanVersion:    ttnpb.MACVersion_MAC_V1_0_2,
				LorawanPhyVersion: ttnpb.PHYVersion_RP001_V1_0_2_REV_B,
				MacSettings: &ttnpb.MACSettings{
					DesiredMaxEirp: &ttnpb.DeviceEIRPValue{
						Value: ttnpb.DeviceEIRP_DEVICE_EIRP_18,
					},
				},
			},
			MACState: func() *ttnpb.MACState {
				macState := MakeDefaultEU868MACState(ttnpb.Class_CLASS_A, ttnpb.MACVersion_MAC_V1_0_2, ttnpb.PHYVersion_RP001_V1_0_2_REV_B)
				macState.DesiredParameters.MaxEirp = 18
				return macState
			}(),
			FrequencyPlanStore: frequencyplans.NewStore(test.FrequencyPlansFetcher),
		},
		{
			Name: "1.0.2/EU868/multicast/class A",
			Device: &ttnpb.EndDevice{
				FrequencyPlanId:   test.EUFrequencyPlanID,
				LorawanVersion:    ttnpb.MACVersion_MAC_V1_0_2,
				LorawanPhyVersion: ttnpb.PHYVersion_RP001_V1_0_2_REV_B,
				MacSettings: &ttnpb.MACSettings{
					DesiredRx1Delay: &ttnpb.RxDelayValue{
						Value: ttnpb.RxDelay_RX_DELAY_13,
					},
				},
				Multicast: true,
			},
			FrequencyPlanStore: frequencyplans.NewStore(test.FrequencyPlansFetcher),
			ErrorAssertion: func(t *testing.T, err error) bool {
				return assertions.New(t).So(err, should.HaveSameErrorDefinitionAs, errClassAMulticast)
			},
		},
		{
			Name: "1.0.2/EU868/multicast/class B",
			Device: &ttnpb.EndDevice{
				FrequencyPlanId:   test.EUFrequencyPlanID,
				LorawanVersion:    ttnpb.MACVersion_MAC_V1_0_2,
				LorawanPhyVersion: ttnpb.PHYVersion_RP001_V1_0_2_REV_B,
				MacSettings: &ttnpb.MACSettings{
					DesiredRx1Delay: &ttnpb.RxDelayValue{
						Value: ttnpb.RxDelay_RX_DELAY_13,
					},
				},
				Multicast:      true,
				SupportsClassB: true,
			},
			MACState: func() *ttnpb.MACState {
				macState := MakeDefaultEU868MACState(ttnpb.Class_CLASS_A, ttnpb.MACVersion_MAC_V1_0_2, ttnpb.PHYVersion_RP001_V1_0_2_REV_B)
				macState.CurrentParameters.Channels = func() (chs []*ttnpb.MACParameters_Channel) {
					for _, ch := range macState.CurrentParameters.Channels {
						chs = append(chs, &ttnpb.MACParameters_Channel{
							DownlinkFrequency: ch.DownlinkFrequency,
						})
					}
					return chs
				}()
				macState.DesiredParameters = macState.CurrentParameters
				macState.DeviceClass = ttnpb.Class_CLASS_B
				return macState
			}(),
			FrequencyPlanStore: frequencyplans.NewStore(test.FrequencyPlansFetcher),
		},
		{
			Name: "1.0.2/EU868/multicast/class C",
			Device: &ttnpb.EndDevice{
				FrequencyPlanId:   test.EUFrequencyPlanID,
				LorawanVersion:    ttnpb.MACVersion_MAC_V1_0_2,
				LorawanPhyVersion: ttnpb.PHYVersion_RP001_V1_0_2_REV_B,
				MacSettings: &ttnpb.MACSettings{
					DesiredRx1Delay: &ttnpb.RxDelayValue{
						Value: ttnpb.RxDelay_RX_DELAY_13,
					},
				},
				Multicast:      true,
				SupportsClassC: true,
			},
			MACState: func() *ttnpb.MACState {
				macState := MakeDefaultEU868MACState(ttnpb.Class_CLASS_A, ttnpb.MACVersion_MAC_V1_0_2, ttnpb.PHYVersion_RP001_V1_0_2_REV_B)
				macState.CurrentParameters.Channels = func() (chs []*ttnpb.MACParameters_Channel) {
					for _, ch := range macState.CurrentParameters.Channels {
						chs = append(chs, &ttnpb.MACParameters_Channel{
							DownlinkFrequency: ch.DownlinkFrequency,
						})
					}
					return chs
				}()
				macState.DesiredParameters = macState.CurrentParameters
				macState.DeviceClass = ttnpb.Class_CLASS_C
				return macState
			}(),
			FrequencyPlanStore: frequencyplans.NewStore(test.FrequencyPlansFetcher),
		},
		{
			Name: "1.0.3/EU868/factory preset frequencies",
			Device: &ttnpb.EndDevice{
				FrequencyPlanId:   test.EUFrequencyPlanID,
				LorawanVersion:    ttnpb.MACVersion_MAC_V1_0_3,
				LorawanPhyVersion: ttnpb.PHYVersion_RP001_V1_0_3_REV_A,
				MacSettings: &ttnpb.MACSettings{
					FactoryPresetFrequencies: []uint64{
						868100000,
						868500000,
						868700000,
						867100000,
					},
				},
			},
			MACState: func() *ttnpb.MACState {
				macState := MakeDefaultEU868MACState(ttnpb.Class_CLASS_A, ttnpb.MACVersion_MAC_V1_0_3, ttnpb.PHYVersion_RP001_V1_0_3_REV_A)
				macState.CurrentParameters.Channels = func() (chs []*ttnpb.MACParameters_Channel) {
					for _, ch := range macState.CurrentParameters.Channels {
						chs = append(chs, &ttnpb.MACParameters_Channel{
							UplinkFrequency:   ch.UplinkFrequency,
							DownlinkFrequency: ch.DownlinkFrequency,
							MinDataRateIndex:  ch.MinDataRateIndex,
							MaxDataRateIndex:  ch.MaxDataRateIndex,
							EnableUplink: func() bool {
								for _, freq := range [...]uint64{
									868100000,
									868500000,
								} {
									if ch.UplinkFrequency == freq {
										return true
									}
								}
								return false
							}(),
						})
					}
					return append(chs,
						&ttnpb.MACParameters_Channel{
							UplinkFrequency:   868700000,
							DownlinkFrequency: 868700000,
							MaxDataRateIndex:  ttnpb.DataRateIndex_DATA_RATE_5,
							EnableUplink:      true,
						},
						&ttnpb.MACParameters_Channel{
							UplinkFrequency:   867100000,
							DownlinkFrequency: 867100000,
							MaxDataRateIndex:  ttnpb.DataRateIndex_DATA_RATE_5,
							EnableUplink:      true,
						},
					)
				}()
				macState.DesiredParameters.Channels = func() (chs []*ttnpb.MACParameters_Channel) {
					for _, ch := range macState.CurrentParameters.Channels {
						chs = append(chs, &ttnpb.MACParameters_Channel{
							UplinkFrequency:   ch.UplinkFrequency,
							DownlinkFrequency: ch.DownlinkFrequency,
							MinDataRateIndex:  ch.MinDataRateIndex,
							MaxDataRateIndex:  ch.MaxDataRateIndex,
							EnableUplink: func() bool {
								for _, freq := range [...]uint64{
									867100000,
									868100000,
									868300000,
									868500000,
								} {
									if ch.UplinkFrequency == freq {
										return true
									}
								}
								return false
							}(),
						})
					}
					return append(chs,
						&ttnpb.MACParameters_Channel{
							UplinkFrequency:   867300000,
							DownlinkFrequency: 867300000,
							MinDataRateIndex:  ttnpb.DataRateIndex_DATA_RATE_0,
							MaxDataRateIndex:  ttnpb.DataRateIndex_DATA_RATE_5,
							EnableUplink:      true,
						},
						&ttnpb.MACParameters_Channel{
							UplinkFrequency:   867500000,
							DownlinkFrequency: 867500000,
							MinDataRateIndex:  ttnpb.DataRateIndex_DATA_RATE_0,
							MaxDataRateIndex:  ttnpb.DataRateIndex_DATA_RATE_5,
							EnableUplink:      true,
						},
						&ttnpb.MACParameters_Channel{
							UplinkFrequency:   867700000,
							DownlinkFrequency: 867700000,
							MinDataRateIndex:  ttnpb.DataRateIndex_DATA_RATE_0,
							MaxDataRateIndex:  ttnpb.DataRateIndex_DATA_RATE_5,
							EnableUplink:      true,
						},
						&ttnpb.MACParameters_Channel{
							UplinkFrequency:   867900000,
							DownlinkFrequency: 867900000,
							MinDataRateIndex:  ttnpb.DataRateIndex_DATA_RATE_0,
							MaxDataRateIndex:  ttnpb.DataRateIndex_DATA_RATE_5,
							EnableUplink:      true,
						},
					)
				}()
				return macState
			}(),
			FrequencyPlanStore: frequencyplans.NewStore(test.FrequencyPlansFetcher),
		},
		{
			Name: "1.1/EU868",
			Device: &ttnpb.EndDevice{
				FrequencyPlanId:   test.EUFrequencyPlanID,
				LorawanVersion:    ttnpb.MACVersion_MAC_V1_1,
				LorawanPhyVersion: ttnpb.PHYVersion_RP001_V1_1_REV_B,
				MacSettings: &ttnpb.MACSettings{
					DesiredRx1Delay: &ttnpb.RxDelayValue{
						Value: ttnpb.RxDelay_RX_DELAY_13,
					},
				},
			},
			MACState: func() *ttnpb.MACState {
				macState := MakeDefaultEU868MACState(ttnpb.Class_CLASS_A, ttnpb.MACVersion_MAC_V1_1, ttnpb.PHYVersion_RP001_V1_1_REV_B)
				macState.DesiredParameters.Rx1Delay = ttnpb.RxDelay_RX_DELAY_13
				return macState
			}(),
			FrequencyPlanStore: frequencyplans.NewStore(test.FrequencyPlansFetcher),
		},
		{
			Name: "1.1/EU868/multicast/class A",
			Device: &ttnpb.EndDevice{
				FrequencyPlanId:   test.EUFrequencyPlanID,
				LorawanVersion:    ttnpb.MACVersion_MAC_V1_1,
				LorawanPhyVersion: ttnpb.PHYVersion_RP001_V1_1_REV_B,
				MacSettings: &ttnpb.MACSettings{
					DesiredRx1Delay: &ttnpb.RxDelayValue{
						Value: ttnpb.RxDelay_RX_DELAY_13,
					},
				},
				Multicast: true,
			},
			FrequencyPlanStore: frequencyplans.NewStore(test.FrequencyPlansFetcher),
			ErrorAssertion: func(t *testing.T, err error) bool {
				return assertions.New(t).So(err, should.HaveSameErrorDefinitionAs, errClassAMulticast)
			},
		},
		{
			Name: "1.1/EU868/multicast/class B",
			Device: &ttnpb.EndDevice{
				FrequencyPlanId:   test.EUFrequencyPlanID,
				LorawanVersion:    ttnpb.MACVersion_MAC_V1_1,
				LorawanPhyVersion: ttnpb.PHYVersion_RP001_V1_1_REV_B,
				MacSettings: &ttnpb.MACSettings{
					DesiredRx1Delay: &ttnpb.RxDelayValue{
						Value: ttnpb.RxDelay_RX_DELAY_13,
					},
				},
				Multicast:      true,
				SupportsClassB: true,
			},
			MACState: func() *ttnpb.MACState {
				macState := MakeDefaultEU868MACState(ttnpb.Class_CLASS_A, ttnpb.MACVersion_MAC_V1_1, ttnpb.PHYVersion_RP001_V1_1_REV_B)
				macState.CurrentParameters.Channels = func() (chs []*ttnpb.MACParameters_Channel) {
					for _, ch := range macState.CurrentParameters.Channels {
						chs = append(chs, &ttnpb.MACParameters_Channel{
							DownlinkFrequency: ch.DownlinkFrequency,
						})
					}
					return chs
				}()
				macState.DesiredParameters = macState.CurrentParameters
				macState.DeviceClass = ttnpb.Class_CLASS_B
				return macState
			}(),
			FrequencyPlanStore: frequencyplans.NewStore(test.FrequencyPlansFetcher),
		},
		{
			Name: "1.1/EU868/multicast/class C",
			Device: &ttnpb.EndDevice{
				FrequencyPlanId:   test.EUFrequencyPlanID,
				LorawanVersion:    ttnpb.MACVersion_MAC_V1_1,
				LorawanPhyVersion: ttnpb.PHYVersion_RP001_V1_1_REV_B,
				MacSettings: &ttnpb.MACSettings{
					DesiredRx1Delay: &ttnpb.RxDelayValue{
						Value: ttnpb.RxDelay_RX_DELAY_13,
					},
				},
				Multicast:      true,
				SupportsClassC: true,
			},
			MACState: func() *ttnpb.MACState {
				macState := MakeDefaultEU868MACState(ttnpb.Class_CLASS_A, ttnpb.MACVersion_MAC_V1_1, ttnpb.PHYVersion_RP001_V1_1_REV_B)
				macState.CurrentParameters.Channels = func() (chs []*ttnpb.MACParameters_Channel) {
					for _, ch := range macState.CurrentParameters.Channels {
						chs = append(chs, &ttnpb.MACParameters_Channel{
							DownlinkFrequency: ch.DownlinkFrequency,
						})
					}
					return chs
				}()
				macState.DesiredParameters = macState.CurrentParameters
				macState.DeviceClass = ttnpb.Class_CLASS_C
				return macState
			}(),
			FrequencyPlanStore: frequencyplans.NewStore(test.FrequencyPlansFetcher),
		},
		{
			Name: "1.0.2/US915_FSB2",
			Device: &ttnpb.EndDevice{
				FrequencyPlanId:   test.USFrequencyPlanID,
				LorawanVersion:    ttnpb.MACVersion_MAC_V1_0_2,
				LorawanPhyVersion: ttnpb.PHYVersion_RP001_V1_0_2_REV_B,
				MacSettings: &ttnpb.MACSettings{
					DesiredRx1Delay: &ttnpb.RxDelayValue{
						Value: ttnpb.RxDelay_RX_DELAY_13,
					},
				},
			},
			MACState: func() *ttnpb.MACState {
				macState := MakeDefaultUS915FSB2MACState(ttnpb.Class_CLASS_A, ttnpb.MACVersion_MAC_V1_0_2, ttnpb.PHYVersion_RP001_V1_0_2_REV_B)
				macState.DesiredParameters.Rx1Delay = ttnpb.RxDelay_RX_DELAY_13
				return macState
			}(),
			FrequencyPlanStore: frequencyplans.NewStore(test.FrequencyPlansFetcher),
		},
		{
			Name: "1.0.3/US915_FSB2/factory preset frequencies",
			Device: &ttnpb.EndDevice{
				FrequencyPlanId:   test.USFrequencyPlanID,
				LorawanVersion:    ttnpb.MACVersion_MAC_V1_0_3,
				LorawanPhyVersion: ttnpb.PHYVersion_RP001_V1_0_3_REV_A,
				MacSettings: &ttnpb.MACSettings{
					FactoryPresetFrequencies: []uint64{
						904300000,
						904700000,
						904900000,
						905100000,
					},
				},
			},
			MACState: func() *ttnpb.MACState {
				macState := MakeDefaultUS915FSB2MACState(ttnpb.Class_CLASS_A, ttnpb.MACVersion_MAC_V1_0_3, ttnpb.PHYVersion_RP001_V1_0_3_REV_A)
				macState.CurrentParameters.Channels = func() (chs []*ttnpb.MACParameters_Channel) {
					for _, ch := range macState.CurrentParameters.Channels {
						chs = append(chs, &ttnpb.MACParameters_Channel{
							UplinkFrequency:   ch.UplinkFrequency,
							DownlinkFrequency: ch.DownlinkFrequency,
							MinDataRateIndex:  ch.MinDataRateIndex,
							MaxDataRateIndex:  ch.MaxDataRateIndex,
							EnableUplink: func() bool {
								for _, freq := range [...]uint64{
									904300000,
									904700000,
									904900000,
									905100000,
								} {
									if ch.UplinkFrequency == freq {
										return true
									}
								}
								return false
							}(),
						})
					}
					return chs
				}()
				macState.DesiredParameters.Channels = func() (chs []*ttnpb.MACParameters_Channel) {
					for _, ch := range macState.CurrentParameters.Channels {
						chs = append(chs, &ttnpb.MACParameters_Channel{
							UplinkFrequency:   ch.UplinkFrequency,
							DownlinkFrequency: ch.DownlinkFrequency,
							MinDataRateIndex:  ch.MinDataRateIndex,
							MaxDataRateIndex:  ch.MaxDataRateIndex,
							EnableUplink: func() bool {
								for _, freq := range [...]uint64{
									903900000,
									904100000,
									904300000,
									904500000,
									904600000,
									904700000,
									904900000,
									905100000,
									905300000,
								} {
									if ch.UplinkFrequency == freq {
										return true
									}
								}
								return false
							}(),
						})
					}
					return chs
				}()
				return macState
			}(),
			FrequencyPlanStore: frequencyplans.NewStore(test.FrequencyPlansFetcher),
		},
		{
			Name: "1.1/US915_FSB2",
			Device: &ttnpb.EndDevice{
				FrequencyPlanId:   test.USFrequencyPlanID,
				LorawanVersion:    ttnpb.MACVersion_MAC_V1_1,
				LorawanPhyVersion: ttnpb.PHYVersion_RP001_V1_1_REV_B,
				MacSettings: &ttnpb.MACSettings{
					DesiredRx1Delay: &ttnpb.RxDelayValue{
						Value: ttnpb.RxDelay_RX_DELAY_13,
					},
				},
			},
			MACState: func() *ttnpb.MACState {
				macState := MakeDefaultUS915FSB2MACState(ttnpb.Class_CLASS_A, ttnpb.MACVersion_MAC_V1_1, ttnpb.PHYVersion_RP001_V1_1_REV_B)
				macState.DesiredParameters.Rx1Delay = ttnpb.RxDelay_RX_DELAY_13
				return macState
			}(),
			FrequencyPlanStore: frequencyplans.NewStore(test.FrequencyPlansFetcher),
		},
	} {
		tc := tc
		test.RunSubtest(t, test.SubtestConfig{
			Name:     tc.Name,
			Parallel: true,
			Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
				pb := ttnpb.Clone(tc.Device)

				macState, err := NewState(pb, tc.FrequencyPlanStore, &ttnpb.MACSettings{})
				if tc.ErrorAssertion != nil {
					a.So(tc.ErrorAssertion(t, err), should.BeTrue)
				} else {
					a.So(err, should.BeNil)
				}
				a.So(macState, should.Resemble, tc.MACState)
				a.So(pb, should.Resemble, tc.Device)
			},
		})
	}
}

func TestBeaconTimeBefore(t *testing.T) {
	for _, tc := range []struct {
		Time     time.Time
		Expected time.Duration
	}{
		{
			Time:     gpstime.Parse(0),
			Expected: 0,
		},
		{
			Time:     gpstime.Parse(time.Nanosecond),
			Expected: 0,
		},
		{
			Time:     gpstime.Parse(BeaconPeriod - time.Second),
			Expected: 0,
		},
		{
			Time:     gpstime.Parse(BeaconPeriod),
			Expected: BeaconPeriod,
		},
		{
			Time:     gpstime.Parse(BeaconPeriod + time.Second),
			Expected: BeaconPeriod,
		},
		{
			Time:     gpstime.Parse(2*BeaconPeriod - time.Second),
			Expected: BeaconPeriod,
		},
		{
			Time:     gpstime.Parse(10 * BeaconPeriod),
			Expected: 10 * BeaconPeriod,
		},
		{
			Time:     gpstime.Parse(10*BeaconPeriod + time.Second),
			Expected: 10 * BeaconPeriod,
		},
	} {
		tc := tc
		test.RunSubtest(t, test.SubtestConfig{
			Name:     tc.Time.String(),
			Parallel: true,
			Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
				a.So(beaconTimeBefore(tc.Time), should.Equal, tc.Expected)
			},
		})
	}
}

func computePingOffset(beaconTime uint32, devAddr types.DevAddr, pingPeriod uint16) uint16 {
	return test.Must(crypto.ComputePingOffset(beaconTime, devAddr, pingPeriod))
}

func TestNextPingSlotAt(t *testing.T) {
	const beaconTime = 10000 * BeaconPeriod
	beaconAt := gpstime.Parse(beaconTime)
	devAddr := types.DevAddr{0x01, 0x34, 0x07, 0x29}

	pingSlotTime := func(pingPeriod uint16, n uint16) time.Time {
		return beaconAt.Add(tBeaconDelay + beaconReserved + time.Duration(computePingOffset(uint32(beaconTime/time.Second), devAddr, pingPeriod)+n*pingPeriod)*pingSlotLen)
	}

	for _, tc := range []struct {
		Name         string
		Device       *ttnpb.EndDevice
		EarliestAt   time.Time
		ExpectedTime time.Time
		ExpectedOk   bool
	}{
		{
			Name:   "no MAC state/no session/no devAddr",
			Device: &ttnpb.EndDevice{},
		},
		{
			Name: "no session/no devAddr",
			Device: &ttnpb.EndDevice{
				MacState: &ttnpb.MACState{},
			},
		},
		{
			Name: "no devAddr",
			Device: &ttnpb.EndDevice{
				MacState: &ttnpb.MACState{},
				Session:  &ttnpb.Session{},
			},
		},
		{
			Name:       "earliestAt:beaconAt;periodicity:0",
			EarliestAt: beaconAt,
			Device: &ttnpb.EndDevice{
				MacState: &ttnpb.MACState{
					PingSlotPeriodicity: &ttnpb.PingSlotPeriodValue{Value: ttnpb.PingSlotPeriod_PING_EVERY_1S},
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr.Bytes(),
				},
			},
			ExpectedTime: pingSlotTime(1<<5, 0),
			ExpectedOk:   true,
		},
		{
			Name:       "earliestAt:beaconAt;periodicity:1",
			EarliestAt: beaconAt,
			Device: &ttnpb.EndDevice{
				MacState: &ttnpb.MACState{
					PingSlotPeriodicity: &ttnpb.PingSlotPeriodValue{Value: ttnpb.PingSlotPeriod_PING_EVERY_2S},
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr.Bytes(),
				},
			},
			ExpectedTime: pingSlotTime(1<<6, 0),
			ExpectedOk:   true,
		},
		{
			Name:       "earliestAt:beaconAt;periodicity:2",
			EarliestAt: beaconAt,
			Device: &ttnpb.EndDevice{
				MacState: &ttnpb.MACState{
					PingSlotPeriodicity: &ttnpb.PingSlotPeriodValue{Value: ttnpb.PingSlotPeriod_PING_EVERY_4S},
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr.Bytes(),
				},
			},
			ExpectedTime: pingSlotTime(1<<7, 0),
			ExpectedOk:   true,
		},
		{
			Name:       "earliestAt:beaconAt;periodicity:3",
			EarliestAt: beaconAt,
			Device: &ttnpb.EndDevice{
				MacState: &ttnpb.MACState{
					PingSlotPeriodicity: &ttnpb.PingSlotPeriodValue{Value: ttnpb.PingSlotPeriod_PING_EVERY_8S},
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr.Bytes(),
				},
			},
			ExpectedTime: pingSlotTime(1<<8, 0),
			ExpectedOk:   true,
		},
		{
			Name:       "earliestAt:beaconAt;periodicity:4",
			EarliestAt: beaconAt,
			Device: &ttnpb.EndDevice{
				MacState: &ttnpb.MACState{
					PingSlotPeriodicity: &ttnpb.PingSlotPeriodValue{Value: ttnpb.PingSlotPeriod_PING_EVERY_16S},
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr.Bytes(),
				},
			},
			ExpectedTime: pingSlotTime(1<<9, 0),
			ExpectedOk:   true,
		},
		{
			Name:       "earliestAt:beaconAt+12s120ms15ns;periodicity:3",
			EarliestAt: beaconAt.Add(12*time.Second + 120*time.Millisecond + 1*time.Microsecond + 500*time.Nanosecond),
			Device: &ttnpb.EndDevice{
				MacState: &ttnpb.MACState{
					PingSlotPeriodicity: &ttnpb.PingSlotPeriodValue{Value: ttnpb.PingSlotPeriod_PING_EVERY_8S},
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr.Bytes(),
				},
			},
			ExpectedTime: pingSlotTime(1<<8, 1),
			ExpectedOk:   true,
		},
		{
			Name:       "earliestAt:beaconAt+20s;periodicity:4",
			EarliestAt: beaconAt.Add(20 * time.Second),
			Device: &ttnpb.EndDevice{
				MacState: &ttnpb.MACState{
					PingSlotPeriodicity: &ttnpb.PingSlotPeriodValue{Value: ttnpb.PingSlotPeriod_PING_EVERY_16S},
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr.Bytes(),
				},
			},
			ExpectedTime: pingSlotTime(1<<9, 1),
			ExpectedOk:   true,
		},
		{
			Name:       "earliestAt:beaconAt+38s;periodicity:4",
			EarliestAt: beaconAt.Add(38 * time.Second),
			Device: &ttnpb.EndDevice{
				MacState: &ttnpb.MACState{
					PingSlotPeriodicity: &ttnpb.PingSlotPeriodValue{Value: ttnpb.PingSlotPeriod_PING_EVERY_16S},
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr.Bytes(),
				},
			},
			ExpectedTime: pingSlotTime(1<<9, 2),
			ExpectedOk:   true,
		},
		{
			Name:       "earliestAt:beaconAt+50s;periodicity:4",
			EarliestAt: beaconAt.Add(50 * time.Second),
			Device: &ttnpb.EndDevice{
				MacState: &ttnpb.MACState{
					PingSlotPeriodicity: &ttnpb.PingSlotPeriodValue{Value: ttnpb.PingSlotPeriod_PING_EVERY_16S},
				},
				Session: &ttnpb.Session{
					DevAddr: devAddr.Bytes(),
				},
			},
			ExpectedTime: pingSlotTime(1<<9, 3),
			ExpectedOk:   true,
		},
	} {
		tc := tc
		test.RunSubtest(t, test.SubtestConfig{
			Name:     tc.Name,
			Parallel: true,
			Func: func(ctx context.Context, t *testing.T, a *assertions.Assertion) {
				ret, ok := NextPingSlotAt(ctx, tc.Device, tc.EarliestAt)
				if !a.So(ok, should.Equal, tc.ExpectedOk) {
					t.FailNow()
				}
				a.So(ret, should.Resemble, tc.ExpectedTime)
				if ok {
					earliestAt := ret
					ret, ok = NextPingSlotAt(ctx, tc.Device, earliestAt)
					a.So(ok, should.BeTrue)
					a.So(ret, should.Resemble, earliestAt)
				}
			},
		})
	}
}
