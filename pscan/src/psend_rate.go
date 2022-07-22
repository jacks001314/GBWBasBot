package pscan

import (
	"time"
)

const (
	MAX_TCP_SYNSCAN_PACKET_SIZE = 58
	DEFAULT_RATE                = 10000
	NSECPERSEC                  = 1000 * 1000 * 1000
)

var psRate *psendRate

type psendRate struct {

	//the number of sending packets thread
	senders uint64

	//the the packets number per second
	rate uint64

	//the bytes to send packets per second
	bandwidth uint64

	//the max size packet to send
	maxPacketSize uint64
}

type PSendRatePerThread struct {
	psRate *psendRate

	sendRate float64

	lastTime uint64

	sleepTime uint64
}

func InitRSendRate(senders, rate, bandwidth uint64) {

	var pktLen uint64 = MAX_TCP_SYNSCAN_PACKET_SIZE

	// Convert specified bandwidth to packet rate. This is an estimate using the
	// max packet size
	if bandwidth > 0 {

		pktLen *= 8

		// 7 byte MAC preamble, 1 byte Start frame, 4 byte CRC, 12 byte
		// inter-frame gap
		pktLen += 8 * 24
		// adjust calculated length if less than the minimum size of an
		// ethernet frame
		if pktLen < 84*8 {
			pktLen = 84 * 8
		}
		// rate is a uint32_t so, don't overflow
		if bandwidth/pktLen > 0xFFFFFFFF {
			rate = 0
		} else {
			rate = bandwidth / pktLen

		}
	}

	if rate == 0 {
		rate = DEFAULT_RATE
	}

	psRate = &psendRate{
		senders:       senders,
		rate:          rate,
		bandwidth:     bandwidth,
		maxPacketSize: MAX_TCP_SYNSCAN_PACKET_SIZE,
	}
}

func now() uint64 {

	return uint64(time.Now().UnixNano())
}

func NewPSendRatePerThread() *PSendRatePerThread {

	var sendRate float64
	var sleepTime uint64 = NSECPERSEC
	var lastTime = now()

	sendRate = float64(psRate.rate) / float64(psRate.senders)

	if psRate.rate > 0 {

		// set the initial time difference
		sleepTime = uint64(NSECPERSEC / sendRate)
		lastTime = now() - uint64(NSECPERSEC/sendRate)

	}

	return &PSendRatePerThread{
		psRate:    psRate,
		sendRate:  sendRate,
		lastTime:  lastTime,
		sleepTime: sleepTime,
	}
}

func (ps *PSendRatePerThread) Sleep() {

	n := now()

	lastRate := float64(1.0 / float64(n-ps.lastTime))

	sleepTime := float64(ps.sleepTime) * ((lastRate / ps.sendRate) + 1) / 2

	time.Sleep(time.Duration(sleepTime) * time.Nanosecond)

	ps.lastTime = n
}
