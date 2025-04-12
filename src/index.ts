import { Bot } from "grammy";
import { config } from "dotenv";
import { scheduleJob } from "node-schedule";
import { botConfig, scheduleConfig } from "./config/config";
import { ReminderHandler } from "./handlers/reminder.handler";

config();

const bot = new Bot(botConfig.token);

const reminderHandler = new ReminderHandler(bot);

// Настройка расписания для напоминаний
scheduleJob(
	{
		hour: scheduleConfig.hour,
		minute: scheduleConfig.minute,
		tz: scheduleConfig.tz,
	},
	() => reminderHandler.sendDebtReminder()
);

// Запуск бота
bot.start();
console.log("Бот запущен!");
console.log("Время до следующей отправки: ", reminderHandler.getTimeUntilNextMessage());
