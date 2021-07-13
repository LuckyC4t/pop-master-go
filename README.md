# pop-master-go

## 前言

最近看到强网杯中有道需要自动化代码审计能力的题——pop-master，正好与自己技能相关，尝试解一波。

## 分析

```php
<?php
include "class.php";
//class.php.txt
highlight_file(__FILE__);
$a = $_GET['pop'];
$b = $_GET['argv'];
$class = unserialize($a);
$class->Dvickz($b);
```
指定函数名调用，去class.php里搜了下，只有一个类中存在这个函数，说明这个类就是入口了。

![image.png](https://i.loli.net/2021/07/13/1i97BAnbOIMTzED.png)

很明显，就是从这16w行代码中，找到一条以xd4a55为起点，拥有eval函数的类为终点的pop链。

### 寻找调用关系

以上图为例，call site的名称是确定的，需要寻找拥有同名函数的类作为`$this->RSphyE8`的值。

可以用嵌套hash表，将类名与类中的方法名作为索引，方法的body作为值。这样只用查询hash表中是否存在对应名称的方法就能跳转到目标方法。

### 污点传播

在本题中，存在将常量直接赋值给变量的情况，这样变量传递到eval中了，因为变量不可控，造成不了危害。所以要进行污点传播，对变量不可控的路径进行剪枝。

在本题中，没有涉及到全局变量，函数中的变量要么来自参数，要么来自类，要么在函数中产生，均不会直接传递到函数外。

所以，每进到一个函数中，就新建一个变量污点状态，以变量名为索引，查询变量的污染情况。因为如果能进到这个函数，说明这个函数的参数被污染了，所以初始化参数的污染情况为true。

#### 污染情况改变

在本题中，变量污染情况改变只来自于赋值，所以对赋值语句进行分析。

赋值语句三种情况：
1. `$left = $right`
2. `$left = $right.'xxxxx'`
3. `$left = 'xxxxx'`

对于3情况，可以直接将`$left`的污染情况置为false; 对于1、2，则可以抽象成当右操作数中存在受污染的变量，则左操作数也受污染。

#### 函数间传播

例如：
```php
    public function wsECXy($dmOxI){
		for($i = 0; $i < 37; $i ++){
			$arNzsS= $dmOxI;
		}
		if(method_exists($this->qhBmwZR, 'mrXAvd')) $this->qhBmwZR->mrXAvd($dmOxI);
		if(method_exists($this->qhBmwZR, 'AovQfA')) $this->qhBmwZR->AovQfA($dmOxI);

    }
```

遇到call site时，判断参数是否受污染，如果受污染，就将目标函数加入到worklist中，反之则无视发生。

如果遇到eval，如果参数受污染，则打印当前调用栈，因为只需找到一条pop链即可，打印完就可以退出。

#### 其他语句

本题中，只涉及到assign、call、eval需要具体分析，所以其他语句可直接利用递归下降的方法来寻找目标语句分析。


### 运行

```bash
go build -o popmaster cmd/popmaster/popmaster.go

./popmaster -file=class.php.txt -class=xd4a55 -method=Dvickz
```

得到：

![image.png](https://i.loli.net/2021/07/13/orpwztVsnHykuG5.png)

与答案一致：

![image.png](https://i.loli.net/2021/07/13/fZ9RN6oaYtQ2yLJ.png)

## 后话

因为这题的限制挺多的，所以适合拿来练手。真实的白盒场景其实挺复杂的，希望有越来越多的人能够做白盒，一起交流一起进步。

## 相关文章

[强网杯[pop_master]与[陀那多]赛题的出题记录](https://www.anquanke.com/post/id/244153)

[第五届强网杯线上赛冠军队 WriteUp - Web 篇](https://mp.weixin.qq.com/s/Y9HdvGtGkr3JCP__pZwWqw)

[pop_master的花式解题思路](https://www.freebuf.com/articles/web/279680.html)