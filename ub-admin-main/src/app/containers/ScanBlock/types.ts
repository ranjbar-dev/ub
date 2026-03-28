/* --- STATE --- */
export interface ScanBlockState {
  isLoading: boolean;
  error: string | null;
}


export interface Currency {
  backgroundImage: string
  code: string
  id: number
  image: string
  mainNetwork: string
  name: string
  otherBlockChainNetworks: {
    code: string
    name: string
  }[]
}

export type ContainerState = ScanBlockState;