Log
=======

注意: v2.5+ 改用接口形式，只要实现 log.LoggerInterface 接口即可，默认已经实现一个基于 logrus 的 logger
.. code:: go

	logger := ytask.Logger.NewYTaskLogger()

	Server := ytask.Server.NewServer(
	    ...
		ytask.Config.Logger(logger),		// 可以不设置 logger
		...
	)


以下为 v2.4 以前版本使用方法：

YTask使用logrus打印日志

输出日志到文件
----------------

.. code:: go

   import (
   "github.com/gojuukaze/YTask/v2/log"
   "github.com/gojuukaze/go-watch-file")

   // write to file
   file,err:=watchFile.OpenWatchFile("xx.log")
   if err != nil {
       panic(err)
   }
   log.YTaskLog.SetOutput(file)

-  `go-watch-file <https://github.com/gojuukaze/go-watch-file>`__
   ：一个专为日志系统编写的读写文件库，会自动监听文件的变化，文件被删除时自动创建新文件。

设置level
----------------

.. code:: go

   // set level
   log.YTaskLog.SetLevel(logrus.InfoLevel)
