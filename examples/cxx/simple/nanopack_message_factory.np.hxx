// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

#ifndef NANOPACK_MESSAGE_FACTORY_HXX
#define NANOPACK_MESSAGE_FACTORY_HXX

#include <memory>
#include <nanopack/message.hxx>
#include <nanopack/reader.hxx>

std::unique_ptr<NanoPack::Message>
make_nanopack_message(NanoPack::Reader &reader);
std::unique_ptr<NanoPack::Message>
make_nanopack_message(NanoPack::Reader &reader, size_t &bytes_read);

#endif
