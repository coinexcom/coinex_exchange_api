#include "coinexapi.h"
#include <QCryptographicHash>
#include <iostream>
#include <QObjectUserData>
#include <QJsonObject>
#include <QJsonDocument>
#include <QSslConfiguration>

CoinexAPI::CoinexAPI(QString accessId, QString secretKey, QObject *parent):
    QObject(parent),
    mAccessID(accessId),
    mSecretKey(secretKey)
{
    setObjectName("CoinexApi");

    mUrl     = "wss://socket.coinex.com";

}

void  CoinexAPI::signin()
{
    QDateTime    currentDateTime = QDateTime::currentDateTimeUtc();
    qint64       mtime           = currentDateTime.toMSecsSinceEpoch();
    auto         sign            = signinGenerate(mtime);
    QVariantMap  param;

    param["id"]     = 1;
    param["method"] = "server.sign";
    param["params"] = QVariantList() << mAccessID << sign << mtime;

    QString  req = variantToJson(param);

    mWebSocket->sendTextMessage(req);
}

void  CoinexAPI::ping()
{
    QVariantMap  param;

    param["id"]     = ++mID;
    param["method"] = "server.ping";
    param["params"] = "[]";

    QString  req = variantToJson(param);

    if (mWebSocket->isValid())
    {
        mWebSocket->sendTextMessage(req);
    }
    else
    {
        qDebug() << "Connection faield!";
    }
}

void  CoinexAPI::serverTime()
{
    QVariantMap  param;

    param["id"]     = ++mID;
    param["method"] = "server.time";
    param["params"] = "[]";

    auto  req = variantToJson(param);

    if (mWebSocket->isValid())
    {
        mWebSocket->sendTextMessage(req);
    }
    else
    {
        qDebug() << "Connection faield!";
    }
}

void  CoinexAPI::subscribe_depth()
{
    QVariantMap  param;

    param["id"]     = ++mID;
    param["method"] = "depth.subscribe";
    param["params"] = QVariantList() << "BTCUSDT" << 10 << "1";

    auto  req = variantToJson(param);

    if (mWebSocket->isValid())
    {
        mWebSocket->sendTextMessage(req);
    }
    else
    {
        qDebug() << "Connection faield!";
    }
}

void  CoinexAPI::subscribe_asset()
{
    qDebug() << "Ping start ";

    QVariantMap  param;

    param["id"]     = ++mID;
    param["method"] = "asset.subscribe";
    param["params"] = QVariantList() << "USDT" << "CET";

    auto  req = variantToJson(param);

    if (mWebSocket->isValid())
    {
        mWebSocket->sendTextMessage(req);
        qDebug() << "Ping start " << mID;
    }
    else
    {
        qDebug() << "Connection faield!";
    }
}

void  CoinexAPI::onConnected()
{
    std::cout << "Connected to coinex ..." << std::endl;

    // Authenticate
    signin();

}

void  CoinexAPI::error(QAbstractSocket::SocketError error)
{
    std::cout << "Error : " << mWebSocket->errorString().toStdString() << std::endl;
}

void  CoinexAPI::closed()
{
    if (mTimer)
    {
        mTimer->stop();
    }

    std::cout << "Connection close !" << std::endl;
}

void  CoinexAPI::onTextMessageReceived(QString message)
{
    qDebug() << message;

    QJsonDocument  res = QJsonDocument::fromJson(message.toLocal8Bit());
    auto           js  = res.object();
    auto           map = res.toVariant().toMap();

    if (js.contains("id") && (js.value("id").toInt() == 1))
    {
        auto  res = js["result"].toObject();

        if (res.value("status").toString() == "success")
        {
            qDebug() << "Authenticate complete. ";

            // subscribe depth
            subscribe_depth();

            // subscribe asset
            subscribe_asset();

            if (mTimer)
            {
                // note: important to keepalive
                mTimer->start(5000);
            }
        }
        else
        {
            qDebug() << "Authenticate not complete : " << js.value("error");
        }
    }
}

QString  CoinexAPI::signinGenerate(qint64 time)
{
    QString  signdata = QString("access_id=%1&tonce=%2&secret_key=%3").arg(mAccessID.toUpper()).arg(time).arg(mSecretKey.toUpper());

    auto  hash = QCryptographicHash::hash(signdata.toLocal8Bit(), QCryptographicHash::Md5).toHex().toUpper();

    return hash;
}

QString  CoinexAPI::variantToJson(QVariantMap map)
{
    QJsonDocument  js = QJsonDocument::fromVariant(map);

    return js.toJson(QJsonDocument::Compact);
}

void  CoinexAPI::start()
{
    mTimer = new QTimer(this);

    connect(mTimer, &QTimer::timeout, this, &CoinexAPI::ping);

    mWebSocket = new QWebSocket;

    connect(mWebSocket, QOverload<QAbstractSocket::SocketError>::of(&QWebSocket::error), this, &CoinexAPI::error);

    connect(mWebSocket, &QWebSocket::connected, this, &CoinexAPI::onConnected);
    connect(mWebSocket, &QWebSocket::disconnected, this, &CoinexAPI::closed);
    connect(mWebSocket, &QWebSocket::textMessageReceived, this, &CoinexAPI::onTextMessageReceived);

    mWebSocket->open(QUrl(mUrl));

// signin();
}
