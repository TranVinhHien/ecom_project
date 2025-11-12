import { OptionType } from "../types/common";

const DropdownValue: OptionType[] = [
    {
        label: 'Trên xuống dưới',
        value: 0,
        name: 'DateCreated',
        code: 'TopToDown',
    },
    {
        label: 'Dưới lên trên',
        value: 1,
        name: 'DateCreated',
        code: 'DownToTop',
    },
]

enum roleE {
    giaovu = "CATECHISM",
    giaovien = "TEACHER",
    truongkhoa = "DEAN",
    truongbomon = "HEAD_OF_DEPARTMENT",
}

enum CheckE {
    ERROR = "ERROR"
}

export { DropdownValue, roleE, CheckE };
