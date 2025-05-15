#include <checker.hpp>
#include <vector>
#include <string>


Work::Work(std::string workType, int id, const std::map<std::string, bool>& workTasks) 
    : type(workType), ID(id), tasksNum(workTasks.size()), tasks(workTasks) {}