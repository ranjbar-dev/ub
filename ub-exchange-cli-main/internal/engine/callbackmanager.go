package engine

import "sync"

type callBackManager struct {
	resultHandler ResultHandler
	mutex         *sync.Mutex
}

func (cbm *callBackManager) callBack(doneOrders []Order, partialOrder *Order) MatchingResult {
	cbm.mutex.Lock()
	defer cbm.mutex.Unlock()
	return cbm.resultHandler.CallBack(doneOrders, partialOrder)
}

func getCallbackManager(resultHandler ResultHandler) *callBackManager {
	return &callBackManager{
		resultHandler: resultHandler,
		mutex:         &sync.Mutex{},
	}
}
