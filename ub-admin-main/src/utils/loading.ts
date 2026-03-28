import { MessageNames, MessageService } from "services/messageService";


export enum LoadingIds{
  ScanButton='scanButton'
}

export const SetLoading=({id,loading}:{id:string|LoadingIds,loading:boolean})=>{
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: id,
    payload: loading,
  });
}