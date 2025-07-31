# Lang - interpreted scripting language
- Simple
- Not-production ready (and ever will be)
- Interpreted and script yeah

# Features
- Functions
- Classes
- Loops (for, while)
- Control (if-else if-else, break)
- Vars
- Syntax
```
import utils > countAllSym;
import test;

func main() {
    var s = "ab1234ation";
    s.capitalize();
    print(s);
    var num = s.substring(1+1, 5);
    num = int(num);
    if (type(num) == "int") {
        print("Yep");
    }

    if (s.contains(string(num))) {
        print("It is in fact contains it");
    } else {
        print("its not");
    }

    print(countAllSym(s, "a"));
}

main();
```


# Plan
- [ ] fix imports
    - currently it just reads imported file to current global env
    - [x] full file import
    - [x] selective import
    - [ ] aliases
    - [ ] proper look ups
- [ ] support big numbers
