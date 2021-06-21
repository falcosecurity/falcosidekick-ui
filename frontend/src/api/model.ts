import { DisplayMode } from '@/types'

export enum Priority {
    NONE = 'none',
    DEBUG = 'debug',
    INFORMATIONAL = 'informational',
    NOTICE = 'notice',
    WARNING = 'warning',
    ERROR = 'error',
    CRITICAL = 'critical',
    ALERT = 'alert',
    EMERGENCY = 'emergency',
}

export type RawEvent = {
    output: string;
    output_fields: null|Array<{ [key: string]: string|null|number }>;
    priority: string;
    rule: string;
    time: string;
}

export type FalcoEvent = {
    output: string;
    outputFields: { [key: string]: string };
    priority: Priority;
    rule: string;
    time: Date;
    filterTime: number;
}

export type Stats = {
    [key in Priority]: number;
}

export type SocketMessage = {
    events?: RawEvent[];
    outputs?: string[];
    retention: number;
    stats?: Stats;
}

export type EventResponse = {
    events?: FalcoEvent[];
    outputs?: string[];
    retention: number;
    stats?: Stats;
}

export type Config = {
    displayMode: DisplayMode;
}
