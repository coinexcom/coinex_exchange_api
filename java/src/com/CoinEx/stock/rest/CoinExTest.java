package com.CoinEx.stock.rest;

import java.io.IOException;
import java.util.logging.Level;
import java.util.logging.Logger;

import org.apache.http.HttpException;

public class CoinExTest {

	public static void main(String[] args) throws HttpException, IOException{
		
	    String access_id = "9578C0CF6F6847C6A6E7D22128FC6BDA";  //CoinEx申请的apiKey
       	String secret_key = "65DC683285FD400FAE8F81860094D69012C3C8EFA0792D7A";  //CoinEx 申请的secret_key
// 	    String url_prex = "https://api.coinex.com/v1/";
		String url_prex = "https://testapi.coinex.com/v1/";
	
	    StockApi api = new StockApi(url_prex, access_id, secret_key);
	    String account = api.account();
	    Logger.getGlobal().log(Level.INFO, account);

	}
}
