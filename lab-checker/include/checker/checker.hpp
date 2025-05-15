#ifndef func_hpp
#define func_hpp

#include <map>
#include <string>

class Work {
private:
    std::string type;
    int ID;
    int tasksNum;
    std::map<std::string, bool> tasks;

public:
    Work(std::string workType, int id, const std::map<std::string, bool>& workTasks);
};

#endif