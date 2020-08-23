# Protobuf

## Defining Your Protocol Format

### Defining Your Protocol Format (proto 2)

To create your address book application, you'll need to start with a .proto file. The definitions in a .proto file are simple: you add a message for each data structure you want to serialize, then specify a name and a type for each field in the message. Here is the .proto file that defines your messages, addressbook.proto

The .proto file starts with a package declaration, which helps to prevent naming conflicts between different projects. In C++, your generated classes will be placed in a namespace matching the package name.

Next, you have your message definitions. A message is just an aggregate containing a set of typed fields. Many standard simple data types are available as field types, including bool, int32, float, double, and string. You can also add further structure to your messages by using other message types as field types – in the above example the Person message contains PhoneNumber messages, while the AddressBook message contains Person messages. You can even define message types nested inside other messages – as you can see, the PhoneNumber type is defined inside Person. You can also define enum types if you want one of your fields to have one of a predefined list of values – here you want to specify that a phone number can be one of MOBILE, HOME, or WORK.

The " = 1", " = 2" markers on each element identify the unique "tag" that field uses in the binary encoding. _Tag numbers 1-15 require one less byte to encode than higher numbers, so as an optimization you can decide to use those tags for the commonly used or repeated elements, leaving tags 16 and higher for less-commonly used optional elements_. Each element in a repeated field requires re-encoding the tag number, so repeated fields are particularly good candidates for this optimization.

Each field must be annotated with one of the following modifiers:

- __required__: a value for the field must be provided, otherwise the message will be considered "uninitialized". If libprotobuf is compiled in debug mode, serializing an uninitialized message will cause an assertion failure. In optimized builds, the check is skipped and the message will be written anyway. However, parsing an uninitialized message will always fail (by returning false from the parse method). Other than this, a required field behaves exactly like an optional field.

- __optional__: the field may or may not be set. If an optional field value isn't set, _a default value is used_. For simple types, you can specify your own default value, as we've done for the phone number type in the example. Otherwise, a system default is used: zero for numeric types, the empty string for strings, false for bools. For embedded messages, the default value is always the "default instance" or "prototype" of the message, which has none of its fields set. Calling the accessor to get the value of an optional (or required) field which has not been explicitly set always returns that field's default value.

- __repeated__: the field may be repeated any number of times (including zero). The order of the repeated values will be preserved in the protocol buffer. Think of repeated fields as dynamically sized arrays.

### Defining Your Protocol Format (proto 3)

Required is no longer supported. To compile we need to use:

## Compiling Your Protocol Buffers

Now that you have a .proto, the next thing you need to do is generate the classes you'll need to read and write AddressBook (and hence Person and PhoneNumber) messages. To do this, you need to run the protocol buffer compiler protoc on your .proto:

Now run the compiler, specifying the source directory (where your application's source code lives – the current directory is used if you don't provide a value), the destination directory (where you want the generated code to go; often the same as $SRC_DIR), and the path to your .proto. In this case, you...:

```ps
protoc --proto_path=$SRC_DIR --cpp_out=$DST_DIR $SRC_DIR/addressbook.proto
```

The actual commands to compile our example in c++, c#, pythong and javascript are:

```ps
protoc --proto_path=. --cpp_out=./protobuf/cc ./addressbook.proto

protoc --proto_path=. --csharp_out=./protobuf/c# ./addressbook.proto

protoc --proto_path=. --python_out=./protobuf/py ./addressbook.proto

protoc --proto_path=. --js_out=./protobuf/js ./addressbook.proto

```

To create it in go we need to first install the go protobuf plugin:

```ps
go get google.golang.org/protobuf/cmd/protoc-gen-go

go install google.golang.org/protobuf/cmd/protoc-gen-go
```

Now we can compile it. Here we are using the version 3 of protobuffer:
```ps
protoc --proto_path=. --go_out=./protobuf/go --experimental_allow_proto3_optional ./addressbook.proto

protoc --proto_path=. --go_out=./src --experimental_allow_proto3_optional ./addressbook.proto
```

## The Protocol Buffer API

Let's look at some of the generated code and see what classes and functions the compiler has created for you. If you look in addressbook.pb.h, you can see that you have a class for each message you specified in addressbook.proto. Looking closer at the Person class, you can see that the compiler has generated accessors for each field. For example, for the name, id, email, and phones fields, you have these methods:

```cc
  // name
  inline bool has_name() const;
  inline void clear_name();
  inline const ::std::string& name() const;
  inline void set_name(const ::std::string& value);
  inline void set_name(const char* value);
  inline ::std::string* mutable_name();

  // id
  inline bool has_id() const;
  inline void clear_id();
  inline int32_t id() const;
  inline void set_id(int32_t value);

  // email
  inline bool has_email() const;
  inline void clear_email();
  inline const ::std::string& email() const;
  inline void set_email(const ::std::string& value);
  inline void set_email(const char* value);
  inline ::std::string* mutable_email();

  // phones
  inline int phones_size() const;
  inline void clear_phones();
  inline const ::google::protobuf::RepeatedPtrField< ::tutorial::Person_PhoneNumber >& phones() const;
  inline ::google::protobuf::RepeatedPtrField< ::tutorial::Person_PhoneNumber >* mutable_phones();
  inline const ::tutorial::Person_PhoneNumber& phones(int index) const;
  inline ::tutorial::Person_PhoneNumber* mutable_phones(int index);
  inline ::tutorial::Person_PhoneNumber* add_phones();
```

As you can see, the getters have exactly the name as the field in lowercase, and the setter methods begin with `set_`. There are also `has_` methods for each singular (required or optional) field which return true if that field has been set. Finally, each field has a `clear_` method that un-sets the field back to its empty state.

While the numeric id field just has the basic accessor set described above, the name and email fields have a couple of extra methods because they're strings – a `mutable_` getter that lets you get a direct pointer to the string, and an _extra setter_. Note that you can call mutable_email() even if email is not already set; it will be initialized to an empty string automatically. If you had a singular message field in this example, it would also have a `mutable_` method but not a `set_` method.

Repeated fields also have some special methods – if you look at the methods for the repeated phones field, you'll see that you can

- check the repeated field's `_size` (in other words, how many phone numbers are associated with this Person).

- get a specified phone number using its index.

- update an existing phone number at the specified index.

- add another phone number to the message which you can then edit (repeated scalar types have an add_ that just lets you pass in the new value).

# GRPC

## Go

In the protofile we are defining both the __message__ and the __service__. To generate the service we need a plugin:

```ps
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

We can generate the protobuf and the client separatedly:

```ps
protoc --proto_path=. --go_out=./src --experimental_allow_proto3_optional ./helloworld.proto
```

```ps
protoc --proto_path=. --go-grpc_out=./src --experimental_allow_proto3_optional ./helloworld.proto
```

Or we can do it in one go:

```ps
protoc --proto_path=. --go_out=./src --go-grpc_out=./src --experimental_allow_proto3_optional ./helloworld.proto
``` 
