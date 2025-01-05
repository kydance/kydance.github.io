#include <string>

class Singleton {
public:
    Singleton(Singleton&) = delete;
    void operator=(const Singleton&) = delete;

    Singleton* GetInstance(const std::string val);

private:
    Singleton(const std::string val):val_(val) {}

    std::string val_;

    static Singleton* instance_;
};

static Singleton* instance_ = nullptr;

Singleton* Singleton::GetInstance(const std::string val) {
    if (instance_ == nullptr) {
        return new Singleton(val);
    }
    return instance_;
}

int main() {

}