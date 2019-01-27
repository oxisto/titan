export class IndustryJobs {
    constructor(public jobs?: Map<string, IndustryJob>,
        public corporationID?: number) { }
}

export class IndustryJob {
    constructor(public status?: string,
        public endDate?: Date,
        public startDate?: Date,
        public completedDate?: Date,
        public pauseDate?: Date
    ) { }
}
