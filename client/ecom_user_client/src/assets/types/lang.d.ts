interface Common {
	where: string;
}

interface ValidationType {
	attribute: string;
	date: string;
	digits: string;
	format: string;
	max: string;
	min: string;
	other: string;
	size: string;
	value: string;
	values: string;
}

interface LanguageType {
	lng: string;
}

export type { Common, ValidationType, LanguageType };
