// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: person.proto

#ifndef GOOGLE_PROTOBUF_INCLUDED_person_2eproto
#define GOOGLE_PROTOBUF_INCLUDED_person_2eproto

#include <limits>
#include <string>

#include <google/protobuf/port_def.inc>
#if PROTOBUF_VERSION < 3021000
#error This file was generated by a newer version of protoc which is
#error incompatible with your Protocol Buffer headers. Please update
#error your headers.
#endif
#if 3021012 < PROTOBUF_MIN_PROTOC_VERSION
#error This file was generated by an older version of protoc which is
#error incompatible with your Protocol Buffer headers. Please
#error regenerate this file with a newer version of protoc.
#endif

#include <google/protobuf/port_undef.inc>
#include <google/protobuf/io/coded_stream.h>
#include <google/protobuf/arena.h>
#include <google/protobuf/arenastring.h>
#include <google/protobuf/generated_message_util.h>
#include <google/protobuf/metadata_lite.h>
#include <google/protobuf/generated_message_reflection.h>
#include <google/protobuf/message.h>
#include <google/protobuf/repeated_field.h>  // IWYU pragma: export
#include <google/protobuf/extension_set.h>  // IWYU pragma: export
#include <google/protobuf/unknown_field_set.h>
// @@protoc_insertion_point(includes)
#include <google/protobuf/port_def.inc>
#define PROTOBUF_INTERNAL_EXPORT_person_2eproto
PROTOBUF_NAMESPACE_OPEN
namespace internal {
class AnyMetadata;
}  // namespace internal
PROTOBUF_NAMESPACE_CLOSE

// Internal implementation detail -- do not use these members.
struct TableStruct_person_2eproto {
  static const uint32_t offsets[];
};
extern const ::PROTOBUF_NAMESPACE_ID::internal::DescriptorTable descriptor_table_person_2eproto;
class ProtoPerson;
struct ProtoPersonDefaultTypeInternal;
extern ProtoPersonDefaultTypeInternal _ProtoPerson_default_instance_;
PROTOBUF_NAMESPACE_OPEN
template<> ::ProtoPerson* Arena::CreateMaybeMessage<::ProtoPerson>(Arena*);
PROTOBUF_NAMESPACE_CLOSE

// ===================================================================

class ProtoPerson final :
    public ::PROTOBUF_NAMESPACE_ID::Message /* @@protoc_insertion_point(class_definition:ProtoPerson) */ {
 public:
  inline ProtoPerson() : ProtoPerson(nullptr) {}
  ~ProtoPerson() override;
  explicit PROTOBUF_CONSTEXPR ProtoPerson(::PROTOBUF_NAMESPACE_ID::internal::ConstantInitialized);

  ProtoPerson(const ProtoPerson& from);
  ProtoPerson(ProtoPerson&& from) noexcept
    : ProtoPerson() {
    *this = ::std::move(from);
  }

  inline ProtoPerson& operator=(const ProtoPerson& from) {
    CopyFrom(from);
    return *this;
  }
  inline ProtoPerson& operator=(ProtoPerson&& from) noexcept {
    if (this == &from) return *this;
    if (GetOwningArena() == from.GetOwningArena()
  #ifdef PROTOBUF_FORCE_COPY_IN_MOVE
        && GetOwningArena() != nullptr
  #endif  // !PROTOBUF_FORCE_COPY_IN_MOVE
    ) {
      InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }

  static const ::PROTOBUF_NAMESPACE_ID::Descriptor* descriptor() {
    return GetDescriptor();
  }
  static const ::PROTOBUF_NAMESPACE_ID::Descriptor* GetDescriptor() {
    return default_instance().GetMetadata().descriptor;
  }
  static const ::PROTOBUF_NAMESPACE_ID::Reflection* GetReflection() {
    return default_instance().GetMetadata().reflection;
  }
  static const ProtoPerson& default_instance() {
    return *internal_default_instance();
  }
  static inline const ProtoPerson* internal_default_instance() {
    return reinterpret_cast<const ProtoPerson*>(
               &_ProtoPerson_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    0;

  friend void swap(ProtoPerson& a, ProtoPerson& b) {
    a.Swap(&b);
  }
  inline void Swap(ProtoPerson* other) {
    if (other == this) return;
  #ifdef PROTOBUF_FORCE_COPY_IN_SWAP
    if (GetOwningArena() != nullptr &&
        GetOwningArena() == other->GetOwningArena()) {
   #else  // PROTOBUF_FORCE_COPY_IN_SWAP
    if (GetOwningArena() == other->GetOwningArena()) {
  #endif  // !PROTOBUF_FORCE_COPY_IN_SWAP
      InternalSwap(other);
    } else {
      ::PROTOBUF_NAMESPACE_ID::internal::GenericSwap(this, other);
    }
  }
  void UnsafeArenaSwap(ProtoPerson* other) {
    if (other == this) return;
    GOOGLE_DCHECK(GetOwningArena() == other->GetOwningArena());
    InternalSwap(other);
  }

  // implements Message ----------------------------------------------

  ProtoPerson* New(::PROTOBUF_NAMESPACE_ID::Arena* arena = nullptr) const final {
    return CreateMaybeMessage<ProtoPerson>(arena);
  }
  using ::PROTOBUF_NAMESPACE_ID::Message::CopyFrom;
  void CopyFrom(const ProtoPerson& from);
  using ::PROTOBUF_NAMESPACE_ID::Message::MergeFrom;
  void MergeFrom( const ProtoPerson& from) {
    ProtoPerson::MergeImpl(*this, from);
  }
  private:
  static void MergeImpl(::PROTOBUF_NAMESPACE_ID::Message& to_msg, const ::PROTOBUF_NAMESPACE_ID::Message& from_msg);
  public:
  PROTOBUF_ATTRIBUTE_REINITIALIZES void Clear() final;
  bool IsInitialized() const final;

  size_t ByteSizeLong() const final;
  const char* _InternalParse(const char* ptr, ::PROTOBUF_NAMESPACE_ID::internal::ParseContext* ctx) final;
  uint8_t* _InternalSerialize(
      uint8_t* target, ::PROTOBUF_NAMESPACE_ID::io::EpsCopyOutputStream* stream) const final;
  int GetCachedSize() const final { return _impl_._cached_size_.Get(); }

  private:
  void SharedCtor(::PROTOBUF_NAMESPACE_ID::Arena* arena, bool is_message_owned);
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(ProtoPerson* other);

  private:
  friend class ::PROTOBUF_NAMESPACE_ID::internal::AnyMetadata;
  static ::PROTOBUF_NAMESPACE_ID::StringPiece FullMessageName() {
    return "ProtoPerson";
  }
  protected:
  explicit ProtoPerson(::PROTOBUF_NAMESPACE_ID::Arena* arena,
                       bool is_message_owned = false);
  public:

  static const ClassData _class_data_;
  const ::PROTOBUF_NAMESPACE_ID::Message::ClassData*GetClassData() const final;

  ::PROTOBUF_NAMESPACE_ID::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  enum : int {
    kFriendsFieldNumber = 5,
    kFirstNameFieldNumber = 1,
    kMiddleNameFieldNumber = 2,
    kLastNameFieldNumber = 3,
    kAgeFieldNumber = 4,
  };
  // repeated .ProtoPerson friends = 5;
  int friends_size() const;
  private:
  int _internal_friends_size() const;
  public:
  void clear_friends();
  ::ProtoPerson* mutable_friends(int index);
  ::PROTOBUF_NAMESPACE_ID::RepeatedPtrField< ::ProtoPerson >*
      mutable_friends();
  private:
  const ::ProtoPerson& _internal_friends(int index) const;
  ::ProtoPerson* _internal_add_friends();
  public:
  const ::ProtoPerson& friends(int index) const;
  ::ProtoPerson* add_friends();
  const ::PROTOBUF_NAMESPACE_ID::RepeatedPtrField< ::ProtoPerson >&
      friends() const;

  // string first_name = 1;
  void clear_first_name();
  const std::string& first_name() const;
  template <typename ArgT0 = const std::string&, typename... ArgT>
  void set_first_name(ArgT0&& arg0, ArgT... args);
  std::string* mutable_first_name();
  PROTOBUF_NODISCARD std::string* release_first_name();
  void set_allocated_first_name(std::string* first_name);
  private:
  const std::string& _internal_first_name() const;
  inline PROTOBUF_ALWAYS_INLINE void _internal_set_first_name(const std::string& value);
  std::string* _internal_mutable_first_name();
  public:

  // string middle_name = 2;
  void clear_middle_name();
  const std::string& middle_name() const;
  template <typename ArgT0 = const std::string&, typename... ArgT>
  void set_middle_name(ArgT0&& arg0, ArgT... args);
  std::string* mutable_middle_name();
  PROTOBUF_NODISCARD std::string* release_middle_name();
  void set_allocated_middle_name(std::string* middle_name);
  private:
  const std::string& _internal_middle_name() const;
  inline PROTOBUF_ALWAYS_INLINE void _internal_set_middle_name(const std::string& value);
  std::string* _internal_mutable_middle_name();
  public:

  // string last_name = 3;
  void clear_last_name();
  const std::string& last_name() const;
  template <typename ArgT0 = const std::string&, typename... ArgT>
  void set_last_name(ArgT0&& arg0, ArgT... args);
  std::string* mutable_last_name();
  PROTOBUF_NODISCARD std::string* release_last_name();
  void set_allocated_last_name(std::string* last_name);
  private:
  const std::string& _internal_last_name() const;
  inline PROTOBUF_ALWAYS_INLINE void _internal_set_last_name(const std::string& value);
  std::string* _internal_mutable_last_name();
  public:

  // int32 age = 4;
  void clear_age();
  int32_t age() const;
  void set_age(int32_t value);
  private:
  int32_t _internal_age() const;
  void _internal_set_age(int32_t value);
  public:

  // @@protoc_insertion_point(class_scope:ProtoPerson)
 private:
  class _Internal;

  template <typename T> friend class ::PROTOBUF_NAMESPACE_ID::Arena::InternalHelper;
  typedef void InternalArenaConstructable_;
  typedef void DestructorSkippable_;
  struct Impl_ {
    ::PROTOBUF_NAMESPACE_ID::RepeatedPtrField< ::ProtoPerson > friends_;
    ::PROTOBUF_NAMESPACE_ID::internal::ArenaStringPtr first_name_;
    ::PROTOBUF_NAMESPACE_ID::internal::ArenaStringPtr middle_name_;
    ::PROTOBUF_NAMESPACE_ID::internal::ArenaStringPtr last_name_;
    int32_t age_;
    mutable ::PROTOBUF_NAMESPACE_ID::internal::CachedSize _cached_size_;
  };
  union { Impl_ _impl_; };
  friend struct ::TableStruct_person_2eproto;
};
// ===================================================================


// ===================================================================

#ifdef __GNUC__
  #pragma GCC diagnostic push
  #pragma GCC diagnostic ignored "-Wstrict-aliasing"
#endif  // __GNUC__
// ProtoPerson

// string first_name = 1;
inline void ProtoPerson::clear_first_name() {
  _impl_.first_name_.ClearToEmpty();
}
inline const std::string& ProtoPerson::first_name() const {
  // @@protoc_insertion_point(field_get:ProtoPerson.first_name)
  return _internal_first_name();
}
template <typename ArgT0, typename... ArgT>
inline PROTOBUF_ALWAYS_INLINE
void ProtoPerson::set_first_name(ArgT0&& arg0, ArgT... args) {
 
 _impl_.first_name_.Set(static_cast<ArgT0 &&>(arg0), args..., GetArenaForAllocation());
  // @@protoc_insertion_point(field_set:ProtoPerson.first_name)
}
inline std::string* ProtoPerson::mutable_first_name() {
  std::string* _s = _internal_mutable_first_name();
  // @@protoc_insertion_point(field_mutable:ProtoPerson.first_name)
  return _s;
}
inline const std::string& ProtoPerson::_internal_first_name() const {
  return _impl_.first_name_.Get();
}
inline void ProtoPerson::_internal_set_first_name(const std::string& value) {
  
  _impl_.first_name_.Set(value, GetArenaForAllocation());
}
inline std::string* ProtoPerson::_internal_mutable_first_name() {
  
  return _impl_.first_name_.Mutable(GetArenaForAllocation());
}
inline std::string* ProtoPerson::release_first_name() {
  // @@protoc_insertion_point(field_release:ProtoPerson.first_name)
  return _impl_.first_name_.Release();
}
inline void ProtoPerson::set_allocated_first_name(std::string* first_name) {
  if (first_name != nullptr) {
    
  } else {
    
  }
  _impl_.first_name_.SetAllocated(first_name, GetArenaForAllocation());
#ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
  if (_impl_.first_name_.IsDefault()) {
    _impl_.first_name_.Set("", GetArenaForAllocation());
  }
#endif // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  // @@protoc_insertion_point(field_set_allocated:ProtoPerson.first_name)
}

// string middle_name = 2;
inline void ProtoPerson::clear_middle_name() {
  _impl_.middle_name_.ClearToEmpty();
}
inline const std::string& ProtoPerson::middle_name() const {
  // @@protoc_insertion_point(field_get:ProtoPerson.middle_name)
  return _internal_middle_name();
}
template <typename ArgT0, typename... ArgT>
inline PROTOBUF_ALWAYS_INLINE
void ProtoPerson::set_middle_name(ArgT0&& arg0, ArgT... args) {
 
 _impl_.middle_name_.Set(static_cast<ArgT0 &&>(arg0), args..., GetArenaForAllocation());
  // @@protoc_insertion_point(field_set:ProtoPerson.middle_name)
}
inline std::string* ProtoPerson::mutable_middle_name() {
  std::string* _s = _internal_mutable_middle_name();
  // @@protoc_insertion_point(field_mutable:ProtoPerson.middle_name)
  return _s;
}
inline const std::string& ProtoPerson::_internal_middle_name() const {
  return _impl_.middle_name_.Get();
}
inline void ProtoPerson::_internal_set_middle_name(const std::string& value) {
  
  _impl_.middle_name_.Set(value, GetArenaForAllocation());
}
inline std::string* ProtoPerson::_internal_mutable_middle_name() {
  
  return _impl_.middle_name_.Mutable(GetArenaForAllocation());
}
inline std::string* ProtoPerson::release_middle_name() {
  // @@protoc_insertion_point(field_release:ProtoPerson.middle_name)
  return _impl_.middle_name_.Release();
}
inline void ProtoPerson::set_allocated_middle_name(std::string* middle_name) {
  if (middle_name != nullptr) {
    
  } else {
    
  }
  _impl_.middle_name_.SetAllocated(middle_name, GetArenaForAllocation());
#ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
  if (_impl_.middle_name_.IsDefault()) {
    _impl_.middle_name_.Set("", GetArenaForAllocation());
  }
#endif // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  // @@protoc_insertion_point(field_set_allocated:ProtoPerson.middle_name)
}

// string last_name = 3;
inline void ProtoPerson::clear_last_name() {
  _impl_.last_name_.ClearToEmpty();
}
inline const std::string& ProtoPerson::last_name() const {
  // @@protoc_insertion_point(field_get:ProtoPerson.last_name)
  return _internal_last_name();
}
template <typename ArgT0, typename... ArgT>
inline PROTOBUF_ALWAYS_INLINE
void ProtoPerson::set_last_name(ArgT0&& arg0, ArgT... args) {
 
 _impl_.last_name_.Set(static_cast<ArgT0 &&>(arg0), args..., GetArenaForAllocation());
  // @@protoc_insertion_point(field_set:ProtoPerson.last_name)
}
inline std::string* ProtoPerson::mutable_last_name() {
  std::string* _s = _internal_mutable_last_name();
  // @@protoc_insertion_point(field_mutable:ProtoPerson.last_name)
  return _s;
}
inline const std::string& ProtoPerson::_internal_last_name() const {
  return _impl_.last_name_.Get();
}
inline void ProtoPerson::_internal_set_last_name(const std::string& value) {
  
  _impl_.last_name_.Set(value, GetArenaForAllocation());
}
inline std::string* ProtoPerson::_internal_mutable_last_name() {
  
  return _impl_.last_name_.Mutable(GetArenaForAllocation());
}
inline std::string* ProtoPerson::release_last_name() {
  // @@protoc_insertion_point(field_release:ProtoPerson.last_name)
  return _impl_.last_name_.Release();
}
inline void ProtoPerson::set_allocated_last_name(std::string* last_name) {
  if (last_name != nullptr) {
    
  } else {
    
  }
  _impl_.last_name_.SetAllocated(last_name, GetArenaForAllocation());
#ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
  if (_impl_.last_name_.IsDefault()) {
    _impl_.last_name_.Set("", GetArenaForAllocation());
  }
#endif // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  // @@protoc_insertion_point(field_set_allocated:ProtoPerson.last_name)
}

// int32 age = 4;
inline void ProtoPerson::clear_age() {
  _impl_.age_ = 0;
}
inline int32_t ProtoPerson::_internal_age() const {
  return _impl_.age_;
}
inline int32_t ProtoPerson::age() const {
  // @@protoc_insertion_point(field_get:ProtoPerson.age)
  return _internal_age();
}
inline void ProtoPerson::_internal_set_age(int32_t value) {
  
  _impl_.age_ = value;
}
inline void ProtoPerson::set_age(int32_t value) {
  _internal_set_age(value);
  // @@protoc_insertion_point(field_set:ProtoPerson.age)
}

// repeated .ProtoPerson friends = 5;
inline int ProtoPerson::_internal_friends_size() const {
  return _impl_.friends_.size();
}
inline int ProtoPerson::friends_size() const {
  return _internal_friends_size();
}
inline void ProtoPerson::clear_friends() {
  _impl_.friends_.Clear();
}
inline ::ProtoPerson* ProtoPerson::mutable_friends(int index) {
  // @@protoc_insertion_point(field_mutable:ProtoPerson.friends)
  return _impl_.friends_.Mutable(index);
}
inline ::PROTOBUF_NAMESPACE_ID::RepeatedPtrField< ::ProtoPerson >*
ProtoPerson::mutable_friends() {
  // @@protoc_insertion_point(field_mutable_list:ProtoPerson.friends)
  return &_impl_.friends_;
}
inline const ::ProtoPerson& ProtoPerson::_internal_friends(int index) const {
  return _impl_.friends_.Get(index);
}
inline const ::ProtoPerson& ProtoPerson::friends(int index) const {
  // @@protoc_insertion_point(field_get:ProtoPerson.friends)
  return _internal_friends(index);
}
inline ::ProtoPerson* ProtoPerson::_internal_add_friends() {
  return _impl_.friends_.Add();
}
inline ::ProtoPerson* ProtoPerson::add_friends() {
  ::ProtoPerson* _add = _internal_add_friends();
  // @@protoc_insertion_point(field_add:ProtoPerson.friends)
  return _add;
}
inline const ::PROTOBUF_NAMESPACE_ID::RepeatedPtrField< ::ProtoPerson >&
ProtoPerson::friends() const {
  // @@protoc_insertion_point(field_list:ProtoPerson.friends)
  return _impl_.friends_;
}

#ifdef __GNUC__
  #pragma GCC diagnostic pop
#endif  // __GNUC__

// @@protoc_insertion_point(namespace_scope)


// @@protoc_insertion_point(global_scope)

#include <google/protobuf/port_undef.inc>
#endif  // GOOGLE_PROTOBUF_INCLUDED_GOOGLE_PROTOBUF_INCLUDED_person_2eproto