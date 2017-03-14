package main

import "time"

type containerRepository struct {
	containersMap map[string]*ContainerTimer

	addChan       chan *ContainerTimer
	cancelChan    chan string
	removeChan    chan string
	listChan      chan chan []*ContainerTimer
}

type ContainerTimer struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	MaxAge    float64 `json:"maxAge"`
	StartedAt int64 `json:"startedAt"`
	timer     *time.Timer
}

func newContainerRepository() *containerRepository {
	containerRepository := containerRepository{
		containersMap:make(map[string]*ContainerTimer),
		addChan       :make(chan *ContainerTimer),
		cancelChan    :make(chan string),
		removeChan    :make(chan string),
		listChan      :make(chan chan []*ContainerTimer),
	}
	containerRepository.start()
	return &containerRepository
}

func (this *containerRepository)start() {
	go func() {

		for {
			select {
			case containerTimer := <-this.addChan:
				this.add(containerTimer)
			case containerId := <-this.cancelChan:
				this.cancel(containerId)
			case containerId := <-this.removeChan:
				this.remove(containerId)
			case returnChan := <-this.listChan:
				this.list(returnChan)
			}
		}
	}()
}

func (this *containerRepository)Add(containerTimer *ContainerTimer) {
	this.addChan <- containerTimer
}

func (this *containerRepository)add(containerTimer *ContainerTimer) {
	this.containersMap[containerTimer.Id] = containerTimer
}

func (this *containerRepository)Cancel(containerId string) {
	this.cancelChan <- containerId
}

func (this *containerRepository)cancel(containerId string) {
	container, containerExists := this.containersMap[containerId]

	if (containerExists) {
		container.timer.Stop()
		delete(this.containersMap, containerId)
	}
}

func (this *containerRepository)Remove(containerId string) {
	this.removeChan <- containerId
}

func (this *containerRepository)remove(containerId string) {
	delete(this.containersMap, containerId)

}

func (this *containerRepository)List() []*ContainerTimer {
	resultChan := make(chan []*ContainerTimer)
	this.listChan <- resultChan
	return <-resultChan
}

func (this *containerRepository)list(returnChan chan []*ContainerTimer) {
	containerList := make([]*ContainerTimer, len(this.containersMap))
	containerCount := 0
	for _, container := range this.containersMap {
		containerList[containerCount] = container
		containerCount++
	}

	returnChan <- containerList
}