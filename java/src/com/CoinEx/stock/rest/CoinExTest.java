package com.CoinEx.stock.rest;

import java.io.IOException;
import java.util.logging.Level;
import java.util.logging.Logger;

import org.apache.http.HttpException;

public class CoinExTest {

	public static void main(String[] args) throws HttpException, IOException{
		
	    String access_id = ""; 
       	String secret_key = ""; 
	    String url_prex = "https://api.coinex.com/v1/"
	
	    StockApi api = new StockApi(url_prex, access_id, secret_key);
	    String account = api.account();
	    Logger.getGlobal().log(Level.INFO, account);

	}
}
