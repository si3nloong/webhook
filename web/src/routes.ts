import WebhookList from "./pages/WebhookList.svelte";
import WebhookDetail from "./pages/WebhookDetail.svelte";

const routes = {
  "/": WebhookList,
  "/webhook/:id": WebhookDetail,
};

export default routes;
