#ifndef COINEXAPI_H
#define COINEXAPI_H

#include <QAbstractSocket>
#include <QThread>
#include <QTimer>
#include <QWebSocket>


class CoinexAPI: public QObject
{
    Q_OBJECT

public:
    CoinexAPI(QString accessId, QString secretKey, QObject *parent = nullptr);

public slots:
    void  start();

    void  onConnected();

    void  error(QAbstractSocket::SocketError error);

    void  closed();

    void  onTextMessageReceived(QString message);

// Coinex API

public slots:
    void     signin();

    void     ping();

    void     serverTime();

    void     subscribe_depth();

    void     subscribe_asset();

private:
    QString  signinGenerate(qint64 time);

    QString  variantToJson(QVariantMap map);

private:
    QWebSocket *mWebSocket = nullptr;
    QUrl        mUrl;
    QString     mAccessID;
    QString     mSecretKey;
    quint64     mID    = 0;
    QTimer     *mTimer = nullptr;
};

#endif // COINEXAPI_H
