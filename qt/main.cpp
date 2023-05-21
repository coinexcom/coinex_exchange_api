#include <QCoreApplication>

#include "coinexapi.h"

int  main(int argc, char *argv[])
{
    QCoreApplication  a(argc, argv);

    CoinexAPI *coinexAPI = new CoinexAPI("ACCESS_TOKEN", "SECRET_TOKEN");
    coinexAPI->start();

    return a.exec();
}
