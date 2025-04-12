import { config } from "dotenv";
import { RecurrenceSpecObjLit } from "node-schedule";

config();

export const botConfig = {
	token: process.env.TG_BOT_TOKEN || "",
	chatId: process.env.TG_CHAT_ID || "",
	logChatId: process.env.TG_CHAT_ID_LOGS || "",
	stickerFileId: process.env.TG_STICKER_FILE_ID || "",
};

export const aiConfig = {
	apiToken: process.env.AI_API_TOKEN || "",
	apiUrl: "https://openrouter.ai/api/v1/chat/completions",
	model: "openrouter/quasar-alpha",
};

export const scheduleConfig: RecurrenceSpecObjLit = {
	hour: 13,
	minute: 45,
	tz: "Europe/Moscow",
};

export const debtors = [
	{ username: process.env.DEBTOR_1 || "" },
	{ username: process.env.DEBTOR_2 || "" },
];
