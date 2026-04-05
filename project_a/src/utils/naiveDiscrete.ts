import { createDiscreteApi } from "naive-ui";

const { message, notification, dialog } = createDiscreteApi(["message", "notification", "dialog"]);

export { message, notification, dialog };
