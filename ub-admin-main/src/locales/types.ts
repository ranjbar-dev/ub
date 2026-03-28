export type ConvertedToFunctionsType<T> = {
	[P in keyof T]: T[P] extends string
	? () => string
	: ConvertedToFunctionsType<T[P]>;
};
export interface FilterArrayElement {
	disabledCols?: string[];
	hiddenCols?: string[];
	dateCols?: string[];
	countryCols?: string[];
	sortableCols?: string[];
	add00ToDate?: boolean;
	dropDownCols?: {
		id: string;
		substituteId?: string;
		options: {
			name: string;
			value: string;
		}[];
	}[];
}
export interface GridTopTab {
	name: string;
	callObject: unknown;
}
export interface Country {
	code: string;
	fullName: string;
	id: number;
	image: string;
	name: string;
}
