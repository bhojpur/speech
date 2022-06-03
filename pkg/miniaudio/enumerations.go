package miniaudio

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Backend type.
type Backend uint32

// Backend enumeration.
const (
	BackendWasapi = iota
	BackendDsound
	BackendWinmm
	BackendCoreaudio
	BackendSndio
	BackendAudio4
	BackendOss
	BackendPulseaudio
	BackendAlsa
	BackendJack
	BackendAaudio
	BackendOpensl
	BackendWebaudio
	BackendNull
)

// DeviceType type.
type DeviceType uint32

// DeviceType enumeration.
const (
	Playback DeviceType = iota + 1
	Capture
	Duplex
	Loopback
)

// ShareMode type.
type ShareMode uint32

// ShareMode enumeration.
const (
	Shared ShareMode = iota
	Exclusive
)

// PerformanceProfile type.
type PerformanceProfile uint32

// PerformanceProfile enumeration.
const (
	LowLatency PerformanceProfile = iota
	Conservative
)

// FormatType type.
type FormatType uint32

// Format enumeration.
const (
	FormatUnknown FormatType = iota
	FormatU8
	FormatS16
	FormatS24
	FormatS32
	FormatF32
)

// ThreadPriority type.
type ThreadPriority int32

// ThreadPriority enumeration.
const (
	ThreadPriorityIdle     ThreadPriority = -5
	ThreadPriorityLowest   ThreadPriority = -4
	ThreadPriorityLow      ThreadPriority = -3
	ThreadPriorityNormal   ThreadPriority = -2
	ThreadPriorityHigh     ThreadPriority = -1
	ThreadPriorityHighest  ThreadPriority = 0
	ThreadPriorityRealtime ThreadPriority = 1

	ThreadPriorityDefault ThreadPriority = 0
)

// ResampleAlgorithm type.
type ResampleAlgorithm uint32

// ResampleAlgorithm enumeration.
const (
	ResampleAlgorithmLinear ResampleAlgorithm = 0
	ResampleAlgorithmSpeex  ResampleAlgorithm = 1
)

// IOSSessionCategory type.
type IOSSessionCategory uint32

// IOSSessionCategory enumeration.
const (
	IOSSessionCategoryDefault       IOSSessionCategory = iota // AVAudioSessionCategoryPlayAndRecord with AVAudioSessionCategoryOptionDefaultToSpeaker.
	IOSSessionCategoryNone                                    // Leave the session category unchanged.
	IOSSessionCategoryAmbient                                 // AVAudioSessionCategoryAmbient
	IOSSessionCategorySoloAmbient                             // AVAudioSessionCategorySoloAmbient
	IOSSessionCategoryPlayback                                // AVAudioSessionCategoryPlayback
	IOSSessionCategoryRecord                                  // AVAudioSessionCategoryRecord
	IOSSessionCategoryPlayAndRecord                           // AVAudioSessionCategoryPlayAndRecord
	IOSSessionCategoryMultiRoute                              // AVAudioSessionCategoryMultiRoute
)

// IOSSessionCategoryOptions type.
type IOSSessionCategoryOptions uint32

// IOSSessionCategoryOptions enumeration.
const (
	IOSSessionCategoryOptionMixWithOthers                        = 0x01 // AVAudioSessionCategoryOptionMixWithOthers
	IOSSessionCategoryOptionDuckOthers                           = 0x02 // AVAudioSessionCategoryOptionDuckOthers
	IOSSessionCategoryOptionAllowBluetooth                       = 0x04 // AVAudioSessionCategoryOptionAllowBluetooth
	IOSSessionCategoryOptionDefaultToSpeaker                     = 0x08 // AVAudioSessionCategoryOptionDefaultToSpeaker
	IOSSessionCategoryOptionInterruptSpokenAudioAndMixWithOthers = 0x11 // AVAudioSessionCategoryOptionInterruptSpokenAudioAndMixWithOthers
	IOSSessionCategoryOptionAllowBluetoothA2dp                   = 0x20 // AVAudioSessionCategoryOptionAllowBluetoothA2DP
	IOSSessionCategoryOptionAllowAirPlay                         = 0x40 // AVAudioSessionCategoryOptionAllowAirPlay
)
