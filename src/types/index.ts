export interface Debtor {
	username: string;
}

export interface BotConfig {
	token: string;
	chatId: string;
	logChatId: string;
	stickerFileId: string;
}

export interface ScheduleConfig {
	sendHour: number;
	sendMinute: number;
}
