import { FacultyType } from '@assets/interface';
import { MouseEvent } from 'react';
import { OptionType } from './common';
import { LanguageType } from './lang';

interface ConfirmModalRefType {
    show?: (event: MouseEvent, data: any, message: string) => void;
}

interface ConfirmModalType {
    agreeLabel?: string;
    cancelLabel?: string;
    acceptLabel?: string;
    rejectLabel?: string;
    onAccept?: (data: any) => void;
    onReject?: () => void;
}

interface SelectObjectModalRefType {
    show: (items: OptionType[], defaultIndex?: number) => void;
}

interface SelectObjectModalType extends LanguageType {
    title: string;
    onConfirm: (data: OptionType) => void;
}

interface SelectFacultyModalRefType {
    show: (selected?: number) => void;
}

interface SelectFacultyModalType extends LanguageType {
    onConfirm: (data: FacultyType | undefined) => void;
}

export type {
    ConfirmModalRefType,
    ConfirmModalType,
    SelectObjectModalRefType,
    SelectObjectModalType,
    SelectFacultyModalRefType,
    SelectFacultyModalType,
};
