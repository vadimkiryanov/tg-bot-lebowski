import { Bot } from "grammy";
import { config } from "dotenv";
import { scheduleJob } from "node-schedule";

config();

// Константы для времени отправки
const SEND_HOUR = 16;
const SEND_MINUTE = 0;

interface Debtor {
	username: string;
}

const bot = new Bot(process.env.TG_BOT_TOKEN || "");

const debtors: Debtor[] = [
	{ username: process.env.DEBTOR_1 || "" },
	{ username: process.env.DEBTOR_2 || "" },
];

// Функция для форматирования времени
function formatTime(hours: number, minutes: number): string {
	const h = hours.toString().padStart(2, "0");
	const m = minutes.toString().padStart(2, "0");
	return `${h}:${m}`;
}

// Функция для расчета оставшегося времени до следующей отправки
function getTimeUntilNextMessage(): string {
	const now = new Date();
	const targetTime = new Date();

	// Устанавливаем время в московском часовом поясе
	targetTime.setUTCHours(SEND_HOUR - 3, SEND_MINUTE, 0, 0); // МСК = UTC+3

	// Если текущее время уже прошло время отправки, добавляем день
	if (now > targetTime) {
		targetTime.setDate(targetTime.getDate() + 1);
	}

	const diff = targetTime.getTime() - now.getTime();
	const hours = Math.floor(diff / (1000 * 60 * 60));
	const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

	return `${formatTime(hours, minutes)}`;
}

async function sendDebtReminder() {
	const chatId = process.env.TG_CHAT_ID;
	const stickerFileId = process.env.TG_STICKER_FILE_ID;

	if (!chatId || !stickerFileId) {
		console.error("Отсутствуют необходимые переменные окружения");
		return;
	}

	for (const debtor of debtors) {
		if (!debtor.username) continue;

		const message = `Как дела, <a href="https://t.me/${debtor.username}">Лебовски</a>?`;

		try {
			await bot.api.sendMessage(chatId, message, {
				parse_mode: "HTML",
				link_preview_options: { is_disabled: true },
			});

			await bot.api.sendSticker(chatId, stickerFileId);
		} catch (error) {
			console.error(`Ошибка при отправке сообщения пользователю ${debtor.username}:`, error);
		}
	}
}

// Функция для отправки сообщения в чат логирования
async function sendLogMessage(message: string) {
	const logChatId = process.env.TG_CHAT_ID_LOGS;
	if (!logChatId) {
		console.error("Не указан ID чата для логирования");
		return;
	}

	try {
		await bot.api.sendMessage(logChatId, message, {
			parse_mode: "Markdown",
			link_preview_options: { is_disabled: true },
		});
	} catch (error) {
		console.error("Ошибка при отправке сообщения в чат логирования:", error);
	}
}

// Запускаем бота
bot.start();

// Настраиваем ежедневную отправку
scheduleJob(
	{
		hour: SEND_HOUR,
		minute: SEND_MINUTE,
		tz: "Europe/Moscow",
	},
	sendDebtReminder
);

// Отправляем сообщение о старте
const timeUntilNext = getTimeUntilNextMessage();
const startMessage =
	`🤖 *Бот запущен!*\n\n` +
	`⏳ До следующей отправки: ${timeUntilNext}\n` +
	`Бот будет отправлять сообщения в чат: https://web.telegram.org/k/#-${process.env.TG_CHAT_ID?.slice(4)}\n` +
	`⏰ Время отправки: ${formatTime(SEND_HOUR, SEND_MINUTE)} по МСК\n`;

sendLogMessage(startMessage);

console.log("Бот запущен!");
