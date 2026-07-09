const DEMO_HOUSEHOLD_ID = 1;
const DEMO_LIST_ID = 1;
const DEMO_MEMBER_ID = 1;

const apiUrlInput = document.querySelector("#api-url");
const itemForm = document.querySelector("#item-form");
const itemNameInput = document.querySelector("#item-name");
const addButton = document.querySelector("#add-button");
const refreshButton = document.querySelector("#refresh-button");
const itemsList = document.querySelector("#items-list");
const emptyState = document.querySelector("#empty-state");
const enableButton = document.querySelector("#enable-button");
const sendButton = document.querySelector("#send-button");
const statusEl = document.querySelector("#status");

const savedApiUrl = localStorage.getItem("groceries-demo-api-url");
if (savedApiUrl) {
  apiUrlInput.value = savedApiUrl;
}

function setStatus(message) {
  statusEl.textContent = message;
}

function getApiUrl() {
  const apiUrl = apiUrlInput.value.trim().replace(/\/$/, "");
  localStorage.setItem("groceries-demo-api-url", apiUrl);
  return apiUrl;
}

function getItemsUrl() {
  return `${getApiUrl()}/households/${DEMO_HOUSEHOLD_ID}/lists/${DEMO_LIST_ID}/items`;
}

function getItemUrl(itemId) {
  return `${getItemsUrl()}/${itemId}`;
}

function getAuthHeaders(extraHeaders = {}) {
  return {
    ...extraHeaders,
    Authorization: `Bearer ${DEMO_MEMBER_ID}`,
  };
}

function normalizeItems(data) {
  if (Array.isArray(data)) {
    return data;
  }

  if (Array.isArray(data?.items)) {
    return data.items;
  }

  return [];
}

function renderItems(items) {
  itemsList.replaceChildren();
  emptyState.hidden = items.length > 0;

  for (const item of items) {
    const itemId = item.id ?? item.ID;
    const itemName = item.name ?? item.Name ?? "Unnamed item";
    const description = item.description ?? item.Description;
    const li = document.createElement("li");
    const content = document.createElement("div");
    content.className = "item-content";
    const nameEl = document.createElement("span");
    nameEl.className = "item-name";
    nameEl.textContent = itemName;
    content.append(nameEl);

    if (description) {
      const descriptionEl = document.createElement("span");
      descriptionEl.className = "item-description";
      descriptionEl.textContent = description;
      content.append(descriptionEl);
    }

    const deleteButton = document.createElement("button");
    deleteButton.className = "delete-button";
    deleteButton.type = "button";
    deleteButton.textContent = "Delete";
    deleteButton.disabled = itemId == null;
    deleteButton.addEventListener("click", () => {
      deleteItem(itemId, itemName, deleteButton);
    });

    li.append(content, deleteButton);
    itemsList.append(li);
  }
}

async function loadItems() {
  refreshButton.disabled = true;
  setStatus("Loading groceries...");

  try {
    const response = await fetch(getItemsUrl(), {
      headers: getAuthHeaders(),
    });
    if (!response.ok) {
      throw new Error(`List request failed with ${response.status}`);
    }

    const data = await response.json();
    const items = normalizeItems(data);
    renderItems(items);
    setStatus(`${items.length} item${items.length === 1 ? "" : "s"} loaded.`);
  } catch (error) {
    renderItems([]);
    setStatus(error.message);
  } finally {
    refreshButton.disabled = false;
  }
}

async function createItem(name) {
  const response = await fetch(getItemsUrl(), {
    method: "POST",
    headers: getAuthHeaders({
      "Content-Type": "application/json",
    }),
    body: JSON.stringify({ name }),
  });

  if (!response.ok) {
    throw new Error(`Add item failed with ${response.status}`);
  }
}

async function deleteItem(itemId, itemName, button) {
  button.disabled = true;
  setStatus(`Deleting ${itemName}...`);

  try {
    const response = await fetch(getItemUrl(itemId), {
      method: "DELETE",
      headers: getAuthHeaders(),
    });

    if (!response.ok) {
      throw new Error(`Delete item failed with ${response.status}`);
    }

    await loadItems();
  } catch (error) {
    setStatus(error.message);
    button.disabled = false;
  }
}

function urlBase64ToUint8Array(base64String) {
  const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding).replace(/-/g, "+").replace(/_/g, "/");
  const rawData = window.atob(base64);
  const outputArray = new Uint8Array(rawData.length);

  for (let i = 0; i < rawData.length; i += 1) {
    outputArray[i] = rawData.charCodeAt(i);
  }

  return outputArray;
}

async function getVapidPublicKey(apiUrl) {
  const response = await fetch(`${apiUrl}/vapid-public-key`);
  if (!response.ok) {
    throw new Error(`VAPID key request failed with ${response.status}`);
  }

  const data = await response.json();
  if (!data.publicKey) {
    throw new Error("Backend response is missing publicKey");
  }

  return data.publicKey;
}

async function enableNotifications() {
  if (!("serviceWorker" in navigator)) {
    throw new Error("This browser does not support service workers.");
  }

  if (!("PushManager" in window)) {
    throw new Error("This browser does not support Web Push.");
  }

  const permission = await Notification.requestPermission();
  if (permission !== "granted") {
    throw new Error(`Notification permission is ${permission}.`);
  }

  const apiUrl = getApiUrl();
  const registration = await navigator.serviceWorker.register("/sw.js");
  const publicKey = await getVapidPublicKey(apiUrl);

  const existingSubscription = await registration.pushManager.getSubscription();
  if (existingSubscription) {
    await existingSubscription.unsubscribe();
  }

  const subscription = await registration.pushManager.subscribe({
    userVisibleOnly: true,
    applicationServerKey: urlBase64ToUint8Array(publicKey),
  });

  const response = await fetch(`${apiUrl}/subscribe`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(subscription),
  });

  if (!response.ok) {
    throw new Error(`Subscription save failed with ${response.status}`);
  }

  sendButton.disabled = false;
  setStatus("Notifications enabled.");
}

async function sendTestPush() {
  const response = await fetch(`${getApiUrl()}/send-test`, {
    method: "POST",
  });

  if (!response.ok) {
    throw new Error(`Test push failed with ${response.status}`);
  }

  setStatus("Test push requested.");
}

apiUrlInput.addEventListener("change", () => {
  getApiUrl();
  loadItems();
});

itemForm.addEventListener("submit", async (event) => {
  event.preventDefault();

  const name = itemNameInput.value.trim();
  if (!name) {
    itemNameInput.focus();
    return;
  }

  addButton.disabled = true;
  setStatus("Adding item...");

  try {
    await createItem(name);
    itemNameInput.value = "";
    await loadItems();
  } catch (error) {
    setStatus(error.message);
  } finally {
    addButton.disabled = false;
    itemNameInput.focus();
  }
});

refreshButton.addEventListener("click", loadItems);

enableButton.addEventListener("click", async () => {
  enableButton.disabled = true;
  setStatus("Enabling notifications...");

  try {
    await enableNotifications();
  } catch (error) {
    enableButton.disabled = false;
    sendButton.disabled = true;
    setStatus(error.message);
  }
});

sendButton.addEventListener("click", async () => {
  sendButton.disabled = true;
  setStatus("Sending test push...");

  try {
    await sendTestPush();
  } catch (error) {
    setStatus(error.message);
  } finally {
    sendButton.disabled = false;
  }
});

loadItems();
