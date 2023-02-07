package consul

import (
	"time"

	"github.com/bosima/ylog"
	"github.com/hashicorp/consul/api"
)

type consulLeaderElection struct {
	serviceId          string
	serviceName        string
	sessionKey         string
	sessionId          string
	leaderKey          string
	lockDelay          int
	sessionTTLDoneChan chan struct{}
	resultHandler      func(result bool)
}

func NewConsulLeaderElection(serviceId string, serviceName string, lockDelaySeconds int, resultHandler func(result bool)) *consulLeaderElection {
	return &consulLeaderElection{
		leaderKey:          serviceName + "/leader",
		lockDelay:          lockDelaySeconds,
		serviceId:          serviceId,
		serviceName:        serviceName,
		sessionId:          "",
		sessionKey:         serviceId + "/leader",
		sessionTTLDoneChan: make(chan struct{}),
		resultHandler:      resultHandler,
	}
}

func (c *consulLeaderElection) Run() error {
	defer close(c.sessionTTLDoneChan)

	waitIndex := uint64(0)
	errorSleepTime := 3
	waitTime := 60 * time.Second
	for {
		q := &api.QueryOptions{
			WaitIndex: waitIndex,
			WaitTime:  waitTime,
		}
		leaderKV, meta, err := consulClient.KV().Get(c.leaderKey, q)
		if err != nil {
			ylog.Error("Block Get Error", err)

			errorSleepTime *= 2
			if errorSleepTime > 300 {
				errorSleepTime = 300
			}
			time.Sleep(time.Duration(errorSleepTime) * time.Second)
			continue
		}
		errorSleepTime = 3

		if leaderKV == nil {
			ylog.Info("leaderKV no exist")
			waitIndex = uint64(0)
			c.wactchHandler(0, nil)
			continue
		}

		waitIndex = meta.LastIndex
		if leaderKV.Session == "" {
			ylog.Info("leaderKV no session")
			c.wactchHandler(waitIndex, leaderKV)
			waitTime = time.Duration(c.lockDelay) * time.Second
		} else {
			waitTime = 60 * time.Second
		}
	}
}

func (c *consulLeaderElection) wactchHandler(index uint64, result interface{}) {
	needAcquire := true
	if result != nil {
		kv := result.(*api.KVPair)
		if kv.Session != "" {
			needAcquire = false
		}
	}
	ylog.Info("needAcquire:", needAcquire)

	if needAcquire {
		// clear session
		if c.sessionId != "" {
			close(c.sessionTTLDoneChan)
			c.sessionTTLDoneChan = make(chan struct{})
			consulClient.Session().Destroy(c.sessionId, nil)
			c.sessionId = ""
		}

		// create session
		sessionId, _, err := consulClient.Session().Create(&api.SessionEntry{
			Name:      c.sessionKey,
			Behavior:  "release",
			TTL:       "30s",
			LockDelay: time.Duration(c.lockDelay) * time.Second,
		}, nil)

		if err != nil {
			ylog.Error("Create Session Error", err)
			return
		}
		c.sessionId = sessionId

		// keep session alive
		go func() {
			consulClient.Session().RenewPeriodic(
				"20s",
				c.sessionId,
				nil,
				c.sessionTTLDoneChan,
			)
		}()

		// acquire leader key
		isLeader, _, err := consulClient.KV().Acquire(&api.KVPair{
			Key:     c.leaderKey,
			Value:   []byte(c.sessionId),
			Session: c.sessionId,
		}, nil)

		if err != nil {
			ylog.Error("Acquire Error", err)
			return
		}

		c.resultHandler(isLeader)
	}
}
