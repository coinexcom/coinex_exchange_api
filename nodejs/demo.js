const ACCESS_ID = ""; // your access id
const SECRET_KEY = ""; // your secret key


function createDictText(params) {
  var keys = Object.keys(params).sort();
  var qs = keys[0] + "=" + params[keys[0]];
  for (var i = 1; i < keys.length; i++) {
    qs += "&" + keys[i] + "=" + params[keys[i]];
  }
  return qs;
}

const crypto = require("crypto");
function createAuthorization(params) {
  var text = createDictText(params) + "&secret_key=" + SECRET_KEY;
  return crypto
    .createHash("md5")
    .update(text)
    .digest("hex")
    .toUpperCase();
}

const Axios = require("axios");
const axios = Axios.create({
  baseURL: "https://api.coinex.com/v1",
  headers: {
    "User-Agent":
      "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36",
    post: {
      "Content-Type": "application/json",
    },
  },
  timeout: 10000,
});


/**
 *  demo
 */

async function getMarketList() {
  const res = await axios.get("/market/list");
  console.log("market list:\n", JSON.stringify(res.data, null, 2));
}

async function getAccountInfo() {
  const params = {
    access_id: ACCESS_ID,
    tonce: Date.now()
  };
  const res = await axios.get("/balance/info", {
    headers: {
      authorization: createAuthorization(params),
    },
    params: params
  });
  console.log("account info:\n", JSON.stringify(res.data, null, 2));
}

async function placeLimitOrder() {
  const data = {
    access_id: ACCESS_ID,
    tonce: Date.now(),
    account_id: 0,
    market: "BTCUSDT",
    type: "buy",
    amount: "0.001",
    price: "10",
  };
  const res = await axios.post("/order/limit", data, {
    headers: {
      authorization: createAuthorization(data),
    },
  });
  console.log("place limit order:\n", JSON.stringify(res.data, null, 2));
}

async function placeMarketOrder() {
  const data = {
    access_id: ACCESS_ID,
    tonce: Date.now(),
    account_id: 0,
    market: "BTCUSDT",
    type: "buy",
    amount: "0.001",
  };
  const res = await axios.post("/order/market", data, {
    headers: {
      authorization: createAuthorization(data),
    },
  });
  console.log("place market order:\n", JSON.stringify(res.data, null, 2));
}

async function getPendingOrder() {
  const params = {
    access_id: ACCESS_ID,
    tonce: Date.now(),
    account_id: 0,
    market: "BTCUSDT",
    page: 1,
    limit: 10
  };
  const res = await axios.get("/order/pending", {
    headers: {
      authorization: createAuthorization(params),
    },
    params: params
  });
  const pendingOrders = res.data.code === 0 ? res.data.data.data : [];
  console.log("pending orders:\n", JSON.stringify(res.data, null, 2));
  return pendingOrders;
}


async function cancelOrder(market, id) {
  const params = {
    access_id: ACCESS_ID,
    tonce: Date.now(),
    account_id: 0,
    market: market,
    id: id,
  };
  const res = await axios.delete("/order/pending", {
    params: params,
    headers: {
      authorization: createAuthorization(params),
    },
  });
  console.log("cancel order:\n", JSON.stringify(res.data, null, 2));
}


async function demo() {
  await getMarketList();

  await getAccountInfo();

  await placeLimitOrder();

  const orders = await getPendingOrder();
  orders.forEach(function(order) {
    cancelOrder(order.market, order.id);
  });
}
demo();

















