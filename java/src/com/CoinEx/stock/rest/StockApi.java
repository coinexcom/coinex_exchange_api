package com.CoinEx.stock.rest;

import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

import org.apache.http.HttpException;


public class StockApi {

	private String secret_key;
	
	private String accessId;
	
	private String url_prex;
	
	public StockApi(String url_prex, String accessId, String secret_key){
		this.accessId = accessId;
		this.secret_key = secret_key;
		this.url_prex = url_prex;
	}
	
	public StockApi(String url_prex){
		this.url_prex = url_prex;
	}
	
	private final String TICKER_URL = "market/ticker";

	private final String DEPTH_URL = "market/depth";
	
	private final String TRADES_URL = "market/deals";
	
	private final String KLINE_URL = "market/kline";

	private final String BALANCE_URL = "balance/";
	
	private final String PENDING_ORDER_URL = "order/pending";
	
	private final String FINISHED_ORDER_URL = "order/finished";

	private final String PUT_LIMIT_URL = "order/limit";

	private final String PUT_MARKET_URL = "order/market";

	private final String CANCEL_ORDER_URL = "order/pending";
	
	private enum HTTP_METHOD {
	    GET, POST, DELETE;
	}
	
	public static enum MARKET {
		MARKET_ETHCNY("ETHCNY"),
		MARKET_ETHBTC("ETHBTC"),
		MARKET_BTCCNY("BTCCNY");
		
	    private String typeName;  
		  
		MARKET(String typeName) {  
            this.typeName = typeName;  
        } 
		
		public String toString() {  
            return this.typeName;  
        } 
	}
	
	public static enum ORDER_TYPE {
		ORDER_TYPE_SELL("sell"),
		ORDER_TYPE_BUY("buy");
		
	    private String typeName;  
		  
		ORDER_TYPE(String typeName) {  
            this.typeName = typeName;  
        } 
		
		public String toString() {  
            return this.typeName;  
        } 
	}
	
	private String doRequest(String url, Map<String, String> paramMap, HTTP_METHOD method) throws HttpException, IOException {
		if (paramMap == null) {
			paramMap = new HashMap<>();
		}
		paramMap.put("access_id", this.accessId);
		paramMap.put("tonce", new Long(System.currentTimeMillis()).toString());
		String authorization = MD5Util.buildMysignV1(paramMap, this.secret_key);
		HttpUtilManager httpUtil = HttpUtilManager.getInstance();
		switch (method) {
		case GET:
			return httpUtil.requestHttpGet(url_prex, url, paramMap, authorization);
			
		case POST:
			return httpUtil.requestHttpPost(url_prex, url, paramMap, authorization);
			
		case DELETE:
			return httpUtil.requestHttpDelete(url_prex, url, paramMap, authorization);

		default:
			return httpUtil.requestHttpGet(url_prex, url, paramMap, authorization);
		}
	}
	
	public String ticker(MARKET market) throws HttpException, IOException {

		HashMap<String, String> param = new HashMap<>();
		param.put("market", market.toString());
	    return doRequest(TICKER_URL, param, HTTP_METHOD.GET);
	}

	public String depth(MARKET market, String merge, Integer limit) throws HttpException, IOException {
		
		HashMap<String, String> param = new HashMap<>();
		param.put("market", market.toString());
		param.put("merge", merge);
		param.put("limit", String.valueOf(limit));
	    return doRequest(DEPTH_URL, param, HTTP_METHOD.GET);
	}
	
	public String trades(MARKET market) throws HttpException, IOException {
		
		HashMap<String, String> param = new HashMap<>();
		param.put("market", market.toString());
	    return doRequest(TRADES_URL, param, HTTP_METHOD.GET);
	}
	
	public String kline(MARKET market, String type) throws HttpException, IOException {
		
		HashMap<String, String> param = new HashMap<>();
		param.put("market", market.toString());
		param.put("type", type);
	    return doRequest(KLINE_URL, param, HTTP_METHOD.GET);
	}

	public String account() throws HttpException, IOException {
		return doRequest(BALANCE_URL, null, HTTP_METHOD.GET);
	}
	
	public String pendingOrder(MARKET market, int page, int account) throws HttpException, IOException {
		HashMap<String, String> param = new HashMap<>();
		param.put("market", market.toString());
		param.put("page", String.valueOf(page));
		param.put("account", String.valueOf(account));
	    return doRequest(PENDING_ORDER_URL, param, HTTP_METHOD.GET);
	}
	
	public String finishedOrder(MARKET market, String page, String account) throws HttpException, IOException {
		HashMap<String, String> param = new HashMap<>();
		param.put("market", market.toString());
		param.put("page", page);
		param.put("account", account);
	    return doRequest(FINISHED_ORDER_URL, param, HTTP_METHOD.GET);
	}
	
	public String putLimitOrder(MARKET market, ORDER_TYPE type, float amount, float price) throws HttpException, IOException {
		HashMap<String, String> param = new HashMap<>();
		param.put("market", market.toString());
		param.put("type", type.toString());
		param.put("amount", String.valueOf(amount));
		param.put("price", String.valueOf(price));
	    return doRequest(PUT_LIMIT_URL, param, HTTP_METHOD.POST);
	}
	
	public String putMarketOrder(MARKET market, ORDER_TYPE type, float amount) throws HttpException, IOException {
		HashMap<String, String> param = new HashMap<>();
		param.put("market", market.toString());
		param.put("type", type.toString());
		param.put("amount", String.valueOf(amount));
	    return doRequest(PUT_MARKET_URL, param, HTTP_METHOD.POST);
	}
	
	public String cancelOrder(MARKET market, String orderID) throws HttpException, IOException {
		HashMap<String, String> param = new HashMap<>();
		param.put("market", market.toString());
		param.put("order_id", orderID);
	    return doRequest(CANCEL_ORDER_URL, param, HTTP_METHOD.DELETE);
	}
}
