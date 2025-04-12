import { Bot } from "grammy";
import { Debtor } from "../types";
import { botConfig, scheduleConfig, debtors } from "../config/config";
import { formatTime } from "../utils/text.utils";

export class ReminderHandler {
	private bot: Bot;
	private debtors: Debtor[];

	constructor(bot: Bot) {
		this.bot = bot;
		this.debtors = debtors || [];
		console.log("Инициализирован ReminderHandler с", this.debtors.length, "должниками");
	}

	async sendDebtReminder() {
		const { chatId, stickerFileId } = botConfig;

		if (!chatId || !stickerFileId) {
			console.error("Отсутствуют необходимые переменные окружения");
			return;
		}

		if (!this.debtors || this.debtors.length === 0) {
			console.error("Нет должников для отправки напоминаний");
			return;
		}

		console.log(`Начинаю отправку напоминаний для ${this.debtors.length} должников`);

		for (const debtor of this.debtors) {
			if (!debtor.username) {
				console.log("Пропускаю должника без username");
				continue;
			}

			const message = `Как дела, <a href="https://t.me/${debtor.username}">Лебовски</a>?`;
			console.log(`Отправляю сообщение пользователю ${debtor.username}`);

			try {
				await this.bot.api.sendMessage(chatId, message, {
					parse_mode: "HTML",
					link_preview_options: { is_disabled: true },
				});

				await this.bot.api.sendSticker(chatId, stickerFileId);
				console.log(`Сообщение успешно отправлено пользователю ${debtor.username}`);
			} catch (error) {
				console.error(
					`Ошибка при отправке сообщения пользователю ${debtor.username}:`,
					error
				);
			}
		}
	}

	getTimeUntilNextMessage(): string {
		const now = new Date();
		const targetTime = new Date();

		// Устанавливаем время в московском часовом поясе
		targetTime.setUTCHours(
			Number(scheduleConfig.hour) - 3,
			Number(scheduleConfig.minute),
			0,
			0
		);

		if (now > targetTime) {
			targetTime.setDate(targetTime.getDate() + 1);
		}

		const diff = targetTime.getTime() - now.getTime();
		const hours = Math.floor(diff / (1000 * 60 * 60));
		const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

		return formatTime(hours, minutes);
	}
}
