.. YTask documentation master file, created by
   sphinx-quickstart on Fri Jul  9 14:15:44 2021.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

Welcome to YTask's documentation!
=================================

| YTask is an asynchronous task queue for handling distributed jobs in
  golang
| golang异步任务/队列 框架

-  `中文文档 <https://doc.ikaze.cn/YTask>`__
   (中文文档更加全面，优先阅读中文文档)
-  `En Doc <https://github.com/gojuukaze/YTask/wiki>`__
-  `Github <https://github.com/gojuukaze/YTask>`__

安装
-----

.. code:: shell

   go get -u github.com/gojuukaze/YTask/v2


特点
-----

-  简单无侵入
-  方便扩展broker，backend
-  支持所有能被序列化为json的类型（如：int，float，数组，结构体，复杂的结构体嵌套等）
-  支持任务重试，延时任务

架构图
-------

.. image:: _static/architecture_diagram.png

.. toctree::
   :maxdepth: 2
   :caption: 使用指南
   :hidden:

   QuickStart
   server
   client
   retry
   delay
   expire
   broker
   backend
   log
   error

Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`
