/* --- STATE --- */
export interface BalancesState {
  balances: Record<string, unknown> | null;
  transferModalBalances: { balances: IWallet[]; type: string } | null;
  balancesHistory: Record<string, unknown> | null;
  isLoading: boolean;
  error: string | null;
}

export type ContainerState=BalancesState;
export enum WalletTypes {
	Hot='hot',
	Cold='cold',
	Internal='internal',
	External='external'
}
export interface IWallet {
	address: string
	code: string
	createdAt: string
	free: string
	locked: string
	name: string
	tag: string
	updatedAt: string
	network?: string
}