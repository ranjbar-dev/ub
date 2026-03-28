import 'package:flutter/material.dart';

import '../../../generated/colors.gen.dart';
import '../../global/autocompleteModel.dart';
import '../custom/alphabeticListView.dart';
import 'UBCircularImage.dart';
import 'UBText.dart';

class AutoCompleteList extends StatefulWidget {
  final List<AutoCompleteItem> itemList;
  final String title;
  final double itemHeight;
  final void Function(AutoCompleteItem item) onItemSelect;
  const AutoCompleteList({
    Key key,
    @required this.itemList,
    @required this.onItemSelect,
    this.title,
    this.itemHeight = 55.0,
  }) : super(key: key);
  @override
  _AutoCompleteListState createState() => _AutoCompleteListState();
}

class _AutoCompleteListState extends State<AutoCompleteList> {
  List<String> strList = [];
  List<Widget> favouriteList = [];
  List<Widget> normalList = [];
  TextEditingController searchController = TextEditingController();

  @override
  void initState() {
    widget.itemList
        .sort((a, b) => a.name.toLowerCase().compareTo(b.name.toLowerCase()));
    filterList();
    searchController.addListener(() {
      filterList();
    });
    super.initState();
  }

  filterList() {
    List<AutoCompleteItem> items = [];
    items.addAll(widget.itemList);
    favouriteList = [];
    normalList = [];
    strList = [];
    if (searchController.text.isNotEmpty) {
      items.retainWhere((item) {
        if (item.searchPhrase != null) {
          return item.searchPhrase.toLowerCase().contains(
                searchController.text.toLowerCase(),
              );
        }
        return item.name.toLowerCase().contains(
              searchController.text.toLowerCase(),
            );
      });
    }
    items.forEach((item) {
      if (item.favourite == true) {
        favouriteList.add(GestureDetector(
          onTap: () => widget.onItemSelect(item),
          child: ListTile(
            leading: Stack(
              children: <Widget>[
                if (item.image != null)
                  UBCircularImage(
                    imageAddress: item.image,
                  ),
                if (item.image != null)
                  Container(
                      height: 40,
                      width: 40,
                      child: Center(
                        child: Icon(
                          Icons.star,
                          color: Colors.blue[700],
                        ),
                      ))
              ],
            ),
            title: UBText(
              text: item.name,
              color: ColorName.white,
            ),
            subtitle: UBText(
              text: item.desc,
              color: ColorName.white,
              size: 10,
            ),
          ),
        ));
      } else {
        normalList.add(
          GestureDetector(
            onTap: () => widget.onItemSelect(item),
            child: ListTile(
              leading: item.image != null
                  ? UBCircularImage(
                      size: 50,
                      imageAddress: item.image,
                    )
                  : null,
              title: Row(
                children: [
                  UBText(
                    text: item.name ?? '',
                    color: ColorName.white,
                    size: 13,
                    weight: FontWeight.w600,
                  ),
                  if (item.inPerentesis != null)
                    const SizedBox(
                      width: 4,
                    ),
                  if (item.inPerentesis != null)
                    UBText(
                      text: "(${item.inPerentesis})",
                      size: 13,
                      color: ColorName.grey80,
                      weight: FontWeight.w600,
                    ),
                ],
              ),
              subtitle: item.desc == null
                  ? null
                  : UBText(
                      text: item.desc,
                      color: ColorName.white,
                      size: 10,
                    ),
            ),
          ),
        );
        strList.add(item.name);
      }
    });

    setState(() {
      // ignore: unnecessary_statements
      strList;
      // ignore: unnecessary_statements
      favouriteList;
      // ignore: unnecessary_statements
      normalList;
      // ignore: unnecessary_statements
      strList;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Container(
          height: 40,
          padding: EdgeInsets.symmetric(horizontal: 8),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              if (widget.title != null)
                UBText(
                  text: widget.title,
                  size: 13,
                )
              else
                const SizedBox(),
              GestureDetector(
                onTap: () => Navigator.pop(context),
                child: Container(
                  width: 24,
                  height: 24,
                  child: const Icon(
                    Icons.close,
                    color: ColorName.greybf,
                  ),
                ),
              )
            ],
          ),
        ),
        Container(
          child: TextFormField(
            autofocus: true,
            controller: searchController,
            style: const TextStyle(
              fontSize: 12,
              color: ColorName.lightText,
            ),
            decoration: const InputDecoration(
              focusedBorder: const OutlineInputBorder(
                borderSide: const BorderSide(
                  color: ColorName.primaryBlue,
                  width: 1.0,
                ),
              ),
              enabledBorder: const OutlineInputBorder(
                borderSide: const BorderSide(
                  color: ColorName.grey36,
                  width: 1.0,
                ),
              ),
              suffixIcon: const Icon(
                Icons.search,
                color: Colors.grey,
              ),
              filled: true,
              fillColor: ColorName.inputBackground,
              contentPadding: const EdgeInsets.only(
                left: 14.0,
                bottom: 8.0,
                top: 8.0,
              ),
              hintText: 'Search',
            ),
          ),
        ),
        Expanded(
          child: AlphabetListScrollView(
            strList: strList,
            highlightTextStyle: const TextStyle(
              color: ColorName.primaryBlue,
            ),
            showPreview: false,
            itemBuilder: (context, index) {
              return normalList[index];
            },
            indexedHeight: (i) {
              return widget.itemHeight;
            },
            keyboardUsage: true,
            headerWidgetList: <AlphabetScrollListHeader>[
              AlphabetScrollListHeader(
                  widgetList: [
                    //Padding(
                    //  padding: const EdgeInsets.only(bottom: 16.0),
                    //  child:
                    //)
                  ],
                  icon: const Icon(Icons.search),
                  indexedHeaderHeight: (index) => 80),
              if (favouriteList.length > 0)
                AlphabetScrollListHeader(
                    widgetList: favouriteList,
                    icon: const Icon(Icons.star),
                    indexedHeaderHeight: (index) {
                      return 70;
                    }),
            ],
          ),
        ),
      ],
    );
  }
}
